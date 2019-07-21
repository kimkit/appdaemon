package common

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/kimkit/iputil"
	"github.com/kimkit/luactl"
	"github.com/kimkit/lualib"
	"github.com/yuin/gopher-lua"
)

var (
	logLib   *LogLib
	httpLib  *lualib.HttpLib
	redisLib *lualib.RedisLib
	dbLib    *lualib.DBLib
)

func initLuaLib() {
	logLib = NewLogLib()
	httpLib = lualib.NewHttpLib(Http)
	redisLib = lualib.NewRedisLib(Redis)
	dbLib = lualib.NewDBLib(DB)
}

func CreateStateHandler() *lua.LState {
	ls := luactl.DefaultCreateStateHandler()
	ls.SetGlobal("printf", ls.NewFunction(Printf))
	ls.SetGlobal("sleep", ls.NewFunction(lualib.Sleep))
	ls.SetGlobal("newcron", ls.NewFunction(lualib.NewCron))
	ls.SetGlobal("uuid", ls.NewFunction(lualib.UUID))
	ls.SetGlobal("md5", ls.NewFunction(lualib.MD5))
	ls.SetGlobal("trim", ls.NewFunction(lualib.Trim))
	ls.SetGlobal("split", ls.NewFunction(lualib.Split))
	ls.SetGlobal("querybuild", ls.NewFunction(lualib.QueryBuild))
	ls.SetGlobal("queryparse", ls.NewFunction(lualib.QueryParse))
	ls.SetGlobal("jsonencode", ls.NewFunction(lualib.JsonEncode))
	ls.SetGlobal("jsondecode", ls.NewFunction(lualib.JsonDecode))
	ls.SetGlobal("getserveriplist", ls.NewFunction(GetServerIpList))
	ls.SetGlobal("getserveriplistbyname", ls.NewFunction(GetServerIpListByName))
	ls.SetGlobal("newredisproxy", ls.NewFunction(NewRedisProxy))
	logLib.RegisterGlobal(ls, "log")
	httpLib.RegisterGlobal(ls, "http")
	redisLib.RegisterGlobal(ls, "redis")
	dbLib.RegisterGlobal(ls, "db")
	for name, script := range Config.LuaScript.Libs {
		if err := ls.DoString(script); err != nil {
			Logger.LogError("common.CreateStateHandler", "load lua script failed: %v (%s)", err, name)
		}
	}
	return ls
}

func output(ls *lua.LState, prefix, suffix string) {
	filename := fmt.Sprintf(
		"%s%c%s.output",
		Config.LogsDir,
		os.PathSeparator,
		ls.GetGlobal("jobname").(lua.LString),
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

func GetServerIpList(ls *lua.LState) int {
	ips := iputil.GetServerIpList()
	var res []interface{}
	for _, ip := range ips {
		res = append(res, ip.String())
	}
	ls.Push(lualib.GoToLua(ls, res))
	return 1
}

func GetServerIpListByName(ls *lua.LState) int {
	ips := iputil.GetServerIpListByName(ls.ToString(1))
	var res []interface{}
	for _, ip := range ips {
		res = append(res, ip.String())
	}
	ls.Push(lualib.GoToLua(ls, res))
	return 1
}

func NewRedisProxy(ls *lua.LState) int {
	password := ""
	if len(Config.Passwords) > 0 {
		password = Config.Passwords[0]
	}
	client := NewRedis(&RedisConfig{
		Addr:     ls.ToString(1),
		Password: password,
	})
	redisLib := lualib.NewRedisLib(map[string]*redis.Client{"#": client})
	ud := ls.NewUserData()
	ud.Value = redisLib
	index := ls.NewTable()
	index.RawSetString("call", ls.NewFunction(redisCall))
	index.RawSetString("pipeline", ls.NewFunction(redisPipeline))
	meta := ls.NewTable()
	meta.RawSetString("__index", index)
	ud.Metatable = meta
	ls.Push(ud)
	return 1
}

func redisCall(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	redisLib := ud.Value.(*lualib.RedisLib)
	var args []lua.LValue
	for i := 2; i <= ls.GetTop(); i++ {
		args = append(args, ls.Get(i))
	}
	ls.Push(ls.NewFunction(redisLib.Call))
	ls.Push(lua.LString("#"))
	for _, arg := range args {
		ls.Push(arg)
	}
	ls.Call(len(args)+1, 2)
	return 2
}

func redisPipeline(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	redisLib := ud.Value.(*lualib.RedisLib)
	var args []lua.LValue
	for i := 2; i <= ls.GetTop(); i++ {
		args = append(args, ls.Get(i))
	}
	ls.Push(ls.NewFunction(redisLib.Pipeline))
	ls.Push(lua.LString("#"))
	for _, arg := range args {
		ls.Push(arg)
	}
	ls.Call(len(args)+1, 2)
	return 2
}
