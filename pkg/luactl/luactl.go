package luactl

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

type LuaScript struct {
	Md5sum string
	Name   string
	Fp     *lua.FunctionProto
	Ls     *lua.LState
	tag    int64
}

type LuaScriptStoreOptions struct {
	CreateStateHandler func() *lua.LState
	CompileHandler     func(string, string) (*lua.FunctionProto, error)
}

type LuaScriptStore struct {
	Map     sync.Map
	options LuaScriptStoreOptions
	locker  sync.Mutex
	tag     int64
}

func DefaultCreateStateHandler() *lua.LState {
	ls := lua.NewState()
	return ls
}

var (
	LuaScriptNameRegexp = regexp.MustCompile(`^[a-zA-Z]\w*$`)
)

func CheckLuaScriptName(name string) error {
	if !LuaScriptNameRegexp.MatchString(name) {
		return fmt.Errorf("luactl.CheckLuaScriptName: lua script name (%s) invalid", name)
	}
	return nil
}

func DefaultCompileHandler(name, script string) (*lua.FunctionProto, error) {
	stmts, err := parse.Parse(bytes.NewBuffer([]byte(script)), name)
	if err != nil {
		return nil, fmt.Errorf("luactl.DefaultCompileHandler: %v (%s)", err, name)
	}
	fp, err := lua.Compile(stmts, name)
	if err != nil {
		return nil, fmt.Errorf("luactl.DefaultCompileHandler: %v (%s)", err, name)
	}
	return fp, nil
}

func opGetOpCode(inst uint32) int {
	return int(inst >> 26)
}

func opGetArgBx(inst uint32) int {
	return int(inst & 0x3ffff)
}

func CheckSetGlobal(fp *lua.FunctionProto) error {
	i := 0
	for idx, inst := range fp.Code {
		op := opGetOpCode(inst)
		if op == lua.OP_CLOSURE {
			if err := CheckSetGlobal(fp.FunctionPrototypes[i]); err != nil {
				return err
			}
			i++
		} else if op == lua.OP_SETGLOBAL {
			bx := opGetArgBx(inst)
			return fmt.Errorf("luactl.CheckSetGlobal: script should not allow to set global: %q at line: %v", fp.Constants[bx], fp.DbgSourcePositions[idx])
		}
	}
	return nil
}

func NoSetGlobalCompileHandler(name, script string) (*lua.FunctionProto, error) {
	fp, err := DefaultCompileHandler(name, script)
	if err != nil {
		return nil, fmt.Errorf("luactl.NoGlobalSetCompileHandler: %v (%s)", err, name)
	}
	if err := CheckSetGlobal(fp); err != nil {
		return nil, fmt.Errorf("luactl.NoGlobalSetCompileHandler: %v (%s)", err, name)
	}
	return fp, nil
}

func NewLuaScriptStore(options LuaScriptStoreOptions) *LuaScriptStore {
	if options.CreateStateHandler == nil {
		options.CreateStateHandler = DefaultCreateStateHandler
	}
	if options.CompileHandler == nil {
		options.CompileHandler = DefaultCompileHandler
	}
	return &LuaScriptStore{
		options: options,
		tag:     time.Now().UnixNano(),
	}
}

func Md5sum(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (lss *LuaScriptStore) Get(name string) *LuaScript {
	if val, ok := lss.Map.Load(name); ok {
		return val.(*LuaScript)
	}
	return nil
}

func (lss *LuaScriptStore) Add(name, script string) error {
	lss.locker.Lock()
	defer lss.locker.Unlock()

	if err := CheckLuaScriptName(name); err != nil {
		return fmt.Errorf("luactl.LuaScriptStore.Add: %v (%s)", err, name)
	}

	md5sum := Md5sum(script)
	if s := lss.Get(name); s != nil && s.Md5sum == md5sum {
		s.tag = lss.tag
		return nil
	}
	fp, err := lss.options.CompileHandler(name, script)
	if err != nil {
		return fmt.Errorf("luactl.LuaScriptStore.Add: %v (%s)", err, name)
	}
	lss.Map.Store(name, &LuaScript{
		Md5sum: md5sum,
		Name:   name,
		Fp:     fp,
		Ls:     nil,
		tag:    lss.tag,
	})
	return nil
}

func (lss *LuaScriptStore) Delete(name string) {
	lss.Map.Delete(name)
}

func (lss *LuaScriptStore) Clean() {
	lss.locker.Lock()
	defer lss.locker.Unlock()

	var names []string
	lss.Map.Range(func(k, v interface{}) bool {
		if v.(*LuaScript).tag != lss.tag {
			names = append(names, k.(string))
		}
		return true
	})
	for _, name := range names {
		lss.Delete(name)
	}
	time.Sleep(time.Millisecond)
	lss.tag = time.Now().UnixNano()
}

func (lss *LuaScriptStore) Run(name string) error {
	s := lss.Get(name)
	if s == nil {
		return fmt.Errorf("luactl.LuaScriptStore.Run: lua script not found (%s)", name)
	}
	if s.Ls == nil {
		s.Ls = lss.options.CreateStateHandler()
	}
	s.Ls.SetGlobal("scriptname", lua.LString(s.Name))
	s.Ls.SetGlobal("scriptmd5sum", lua.LString(s.Md5sum))
	s.Ls.Push(s.Ls.NewFunctionFromProto(s.Fp))
	if err := s.Ls.PCall(0, 0, nil); err != nil {
		return fmt.Errorf("luactl.LuaScriptStore.Run: %v (%s)", err, name)
	}
	return nil
}
