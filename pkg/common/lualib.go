package common

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/kimkit/luactl"
	"github.com/kimkit/lualib"
	"github.com/yuin/gopher-lua"
)

func CreateStateHandler() *lua.LState {
	ls := luactl.DefaultCreateStateHandler()
	ls.SetGlobal("printf", ls.NewFunction(Printf))
	ls.SetGlobal("sleep", ls.NewFunction(lualib.Sleep))
	ls.SetGlobal("newcron", ls.NewFunction(lualib.NewCron))
	ls.SetGlobal("uuid", ls.NewFunction(lualib.UUID))
	ls.SetGlobal("md5", ls.NewFunction(lualib.MD5))
	ls.SetGlobal("trim", ls.NewFunction(lualib.Trim))
	NewLogLib().RegisterGlobal(ls, "log")
	lualib.NewHttpLib(map[string]*http.Client{"#": HttpClient}).RegisterGlobal(ls, "http")
	lualib.NewRedisLib(map[string]*redis.Client{"#": RedisClient}).RegisterGlobal(ls, "redis")
	return ls
}

func output(ls *lua.LState, prefix, suffix string) {
	filename := fmt.Sprintf(
		"%s%c%s.output",
		Config.LogsDir,
		os.PathSeparator,
		"luascript_"+ls.GetGlobal("scriptname").(lua.LString),
	)
	format := ls.ToString(1)
	top := ls.GetTop()
	v := make([]interface{}, top-1)
	for i := 2; i <= top; i++ {
		v[i-2] = lualib.LuaToGo(ls.Get(i))
	}
	str := prefix + fmt.Sprintf(format, v...) + suffix
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		Logger.LogError("common.output", "%v", err)
		return
	}
	defer fp.Close()
	if _, err := fp.Write([]byte(str)); err != nil {
		Logger.LogError("common.output", "%v", err)
		return
	}
}

func Printf(ls *lua.LState) int {
	output(ls, "", "")
	return 0
}

type LogLib struct{}

func NewLogLib() *LogLib {
	return &LogLib{}
}

func (ll *LogLib) RegisterGlobal(ls *lua.LState, name string) {
	t := ls.NewTable()
	t.RawSetString("info", ls.NewFunction(ll.Info))
	t.RawSetString("error", ls.NewFunction(ll.Error))
	t.RawSetString("warning", ls.NewFunction(ll.Warning))
	t.RawSetString("debug", ls.NewFunction(ll.Debug))
	ls.SetGlobal(name, t)
}

func (ll *LogLib) Info(ls *lua.LState) int {
	output(ls, fmt.Sprintf("%s INFO ", time.Now().Format("2006/01/02 15:04:05")), "\n")
	return 0
}

func (ll *LogLib) Error(ls *lua.LState) int {
	output(ls, fmt.Sprintf("%s ERROR ", time.Now().Format("2006/01/02 15:04:05")), "\n")
	return 0
}

func (ll *LogLib) Warning(ls *lua.LState) int {
	output(ls, fmt.Sprintf("%s WARNING ", time.Now().Format("2006/01/02 15:04:05")), "\n")
	return 0
}

func (ll *LogLib) Debug(ls *lua.LState) int {
	output(ls, fmt.Sprintf("%s DEBUG ", time.Now().Format("2006/01/02 15:04:05")), "\n")
	return 0
}
