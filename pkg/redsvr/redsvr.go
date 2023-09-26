package redsvr

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rokumoe/redisgo"
)

type Conn struct {
	net.Conn
	Map *sync.Map
}

func Log(level, format string, v ...interface{}) {
	log.Print(level, " ", fmt.Sprintf(format, v...))
}

func LogFatal(level, format string, v ...interface{}) {
	log.Print(level, " ", fmt.Sprintf(format, v...))
	os.Exit(1)
}

func writeResp(conn *Conn, resp *redisgo.Resp) error {
	if err := redisgo.EncodeResp(conn, resp); err != nil {
		conn.Close()
		return err
	}
	return nil
}

func WriteSimpleString(conn *Conn, format string, v ...interface{}) error {
	return writeResp(conn, &redisgo.Resp{Kind: redisgo.SimpleKind, Data: fmt.Sprintf(format, v...)})
}

func WriteError(conn *Conn, format string, v ...interface{}) error {
	return writeResp(conn, &redisgo.Resp{Kind: redisgo.ErrorKind, Data: fmt.Sprintf(format, v...)})
}

func WriteInt(conn *Conn, data int) error {
	return writeResp(conn, &redisgo.Resp{Kind: redisgo.IntegerKind, Data: fmt.Sprintf("%d", data)})
}

func WriteBlukString(conn *Conn, format string, v ...interface{}) error {
	return writeResp(conn, &redisgo.Resp{Kind: redisgo.BlukKind, Data: fmt.Sprintf(format, v...)})
}

func buildArray(data interface{}) redisgo.Resp {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Slice {
		var array []redisgo.Resp
		for i := 0; i < dataValue.Len(); i++ {
			array = append(array, buildArray(dataValue.Index(i).Interface()))
		}
		return redisgo.Resp{Kind: redisgo.ArrayKind, Array: array}
	}
	if dataValue.Kind() == reflect.Map {
		var array []redisgo.Resp
		for _, keyValue := range dataValue.MapKeys() {
			array = append(array, buildArray(keyValue.Interface()))
			array = append(array, buildArray(dataValue.MapIndex(keyValue).Interface()))
		}
		return redisgo.Resp{Kind: redisgo.ArrayKind, Array: array}
	}

	switch data.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return redisgo.Resp{Kind: redisgo.IntegerKind, Data: fmt.Sprint(data)}
	default:
		return redisgo.Resp{Kind: redisgo.BlukKind, Data: fmt.Sprint(data)}
	}
}

func WriteArray(conn *Conn, data interface{}) error {
	resp := buildArray(data)
	if resp.Kind != redisgo.ArrayKind {
		resp = redisgo.Resp{Kind: redisgo.ArrayKind, Array: []redisgo.Resp{resp}}
	}
	return writeResp(conn, &resp)
}

func WriteNull(conn *Conn) error {
	return writeResp(conn, &redisgo.Resp{Kind: redisgo.BlukKind, Null: true})
}

type CommandHandler func(*Command, []string, *Conn) error

type Command struct {
	Name      string
	Argc      int
	S1Handler CommandHandler
	S2Handler CommandHandler
}

type CommandInstance struct {
	Command *Command
	Args    []string
	Conn    *Conn
}

type Server struct {
	commands             map[string]*Command
	commandInstanceQueue chan *CommandInstance
	exitChan             chan struct{}
	closeConnChan        chan struct{}
	pkgSize              int
	wg                   sync.WaitGroup
}

func NewServer() *Server {
	svr := &Server{
		commands:             make(map[string]*Command),
		commandInstanceQueue: make(chan *CommandInstance),
		exitChan:             make(chan struct{}),
		closeConnChan:        make(chan struct{}),
		pkgSize:              10240,
	}
	return svr
}

func (svr *Server) SetPkgSize(size int) {
	svr.pkgSize = size
}

func (svr *Server) Register(cmd *Command) {
	svr.commands[strings.ToLower(cmd.Name)] = cmd
}

func (svr *Server) Unregister(name string) {
	delete(svr.commands, strings.ToLower(name))
}

func (svr *Server) s2() {
	defer svr.wg.Done()
	for {
		select {
		case <-svr.exitChan:
			close(svr.closeConnChan)
			return
		case ci := <-svr.commandInstanceQueue:
			if ci.Command != nil {
				if ci.Command.S2Handler != nil {
					func() {
						defer func() {
							if err := recover(); err != nil {
								log.Printf("[PANIC] %v", err)
								debug.PrintStack()
								ci.Conn.Close()
							}
						}()
						if err := ci.Command.S2Handler(ci.Command, ci.Args, ci.Conn); err != nil {
							WriteError(ci.Conn, "ERR %v", err)
						}
					}()
				}
			}
		}
	}
}

func (svr *Server) s1watcher(closeConnChan chan struct{}, conn *Conn) {
	defer svr.wg.Done()
	select {
	case <-svr.closeConnChan:
		conn.Close()
	case <-closeConnChan:
		// pass
	}
}

func (svr *Server) s1(closeConnChan chan struct{}, conn *Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[PANIC] %v", err)
			debug.PrintStack()
			conn.Close()
		}
		close(closeConnChan)
		svr.wg.Done()
	}()
	decoder := redisgo.NewDecoder(conn, svr.pkgSize)
	for {
		var resp redisgo.Resp
		if err := decoder.Decode(&resp); err != nil {
			conn.Close()
			return
		}
		if len(resp.Array) == 0 {
			conn.Close()
			return
		}
		var name string
		var args []string
		for i, r := range resp.Array {
			if r.Kind != redisgo.BlukKind {
				conn.Close()
				return
			}
			if i == 0 {
				name = strings.ToLower(r.Data)
			} else {
				args = append(args, r.Data)
			}
		}
		if _, ok := svr.commands[name]; !ok {
			WriteError(conn, "ERR unknown command '%s'", name)
			continue
		}
		if svr.commands[name].Argc >= 0 {
			if svr.commands[name].Argc != len(args) {
				WriteError(conn, "ERR wrong number of arguments for '%s'", name)
				continue
			}
		}
		if svr.commands[name].S1Handler != nil {
			if err := svr.commands[name].S1Handler(svr.commands[name], args, conn); err != nil {
				WriteError(conn, "ERR %v", err)
				continue
			}
		}
		if svr.commands[name].S2Handler == nil {
			continue
		}
		ci := &CommandInstance{svr.commands[name], args, conn}
		select {
		case svr.commandInstanceQueue <- ci:
			// pass
		case <-time.After(time.Second):
			WriteError(conn, "ERR timeout for '%s'", name)
		}
	}
}

func (svr *Server) accept(ln net.Listener) {
	defer svr.wg.Done()
	defer close(svr.exitChan)

	svr.wg.Add(1)
	go svr.s2()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		connWrapper := &Conn{conn, &sync.Map{}}
		closeConnChan := make(chan struct{})
		svr.wg.Add(1)
		go svr.s1watcher(closeConnChan, connWrapper)
		svr.wg.Add(1)
		go svr.s1(closeConnChan, connWrapper)
	}
}

func (svr *Server) Accept(ln net.Listener) {
	svr.wg.Add(1)
	go svr.accept(ln)
}

func (svr *Server) Wait() {
	svr.wg.Wait()
}

func (svr *Server) ListenAndServe(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		LogFatal("ERROR", "%v", err)
	}

	svr.Accept(ln)

	Log("INFO", "server started (listening on %s)", addr)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-sigChan
	if sig == os.Interrupt {
		fmt.Println("")
	}

	Log("INFO", "server shutting down ...")

	ln.Close()
	svr.Wait()

	Log("INFO", "server stopped")
}

func NewCommand(ptr interface{}) *Command {
	cmd := &Command{}
	cmdVal := reflect.ValueOf(cmd)
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Kind() != reflect.Ptr {
		return nil
	}
	if ptrVal.Elem().Kind() != reflect.Struct {
		return nil
	}
	if nameVal := ptrVal.Elem().FieldByName("Name"); nameVal.Kind() == reflect.String {
		cmdVal.Elem().FieldByName("Name").Set(nameVal)
	} else {
		cmd.Name = strings.TrimSuffix(strings.ToLower(ptrVal.Elem().Type().Name()), "command")
	}
	if argcVal := ptrVal.Elem().FieldByName("Argc"); argcVal.Kind() == reflect.Int {
		cmdVal.Elem().FieldByName("Argc").Set(argcVal)
	} else {
		cmd.Argc = -1
	}
	s1HandlerVal := ptrVal.MethodByName("S1Handler")
	if s1HandlerVal.IsValid() && s1HandlerVal.Type().AssignableTo(reflect.TypeOf(cmd.S1Handler)) {
		cmdVal.Elem().FieldByName("S1Handler").Set(s1HandlerVal)
	}
	s2HandlerVal := ptrVal.MethodByName("S2Handler")
	if s2HandlerVal.IsValid() && s2HandlerVal.Type().AssignableTo(reflect.TypeOf(cmd.S2Handler)) {
		cmdVal.Elem().FieldByName("S2Handler").Set(s2HandlerVal)
	}
	return cmd
}
