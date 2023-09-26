package ginsvr

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
}

func init() {
	log.SetFlags(log.LstdFlags)
	gin.SetMode(gin.ReleaseMode)
}

func NewServer() *Server {
	svr := &Server{Engine: gin.New()}
	svr.Engine.Use(gin.Recovery())
	return svr
}

func (svr *Server) Register(ptr interface{}) {
	var handler func(*gin.Context)
	handlerType := reflect.TypeOf(handler)
	engineVal := reflect.ValueOf(svr.Engine)
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Kind() != reflect.Ptr {
		return
	}
	if ptrVal.Elem().Kind() != reflect.Struct {
		return
	}
	pathVal := ptrVal.Elem().FieldByName("Path")
	if pathVal.Kind() != reflect.String {
		return
	}
	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"} {
		methodVal := ptrVal.MethodByName(method)
		if methodVal.IsValid() && methodVal.Type().AssignableTo(handlerType) {
			engineVal.MethodByName(method).Call([]reflect.Value{pathVal, methodVal})
		}
	}
}

func (svr *Server) ListenAndServe(addr string) {
	log.Printf("INFO server started (listening on %s)", addr)
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
		sig := <-sigChan
		if sig == os.Interrupt {
			fmt.Println("")
		}
	}()
	if err := gracehttp.Serve(&http.Server{Addr: addr, Handler: svr.Engine}); err != nil {
		log.Fatalf("ERROR %v", err)
	}
	log.Printf("INFO server stopped")
}
