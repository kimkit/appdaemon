package common

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/kimkit/luactl"
	"github.com/yuin/gopher-lua"
)

func CreateStateHandler() *lua.LState {
	ls := luactl.DefaultCreateStateHandler()
	ls.SetGlobal("sleep", ls.NewFunction(LuaLibSleep))
	ls.SetGlobal("printf", ls.NewFunction(LuaLibPrintf))
	ls.SetGlobal("newcron", ls.NewFunction(LuaLibNewCron))
	return ls
}

func LuaToGo(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		if v == nil {
			return nil
		}
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]interface{})
			v.ForEach(func(key, value lua.LValue) {
				ret[fmt.Sprint(LuaToGo(key))] = LuaToGo(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, LuaToGo(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}

func GoToLua(ls *lua.LState, gv interface{}) lua.LValue {
	switch v := gv.(type) {
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case json.Number:
		return lua.LString(v)
	case []interface{}:
		t := ls.NewTable()
		for _, vv := range v {
			t.Append(GoToLua(ls, vv))
		}
		return t
	case map[string]interface{}:
		t := ls.NewTable()
		for vk, vv := range v {
			t.RawSetString(vk, GoToLua(ls, vv))
		}
		return t
	case nil:
		return lua.LNil
	default:
		return lua.LNil
	}
}

func LuaLibSleep(ls *lua.LState) int {
	t := ls.ToInt(1)
	if t < 0 {
		t = 0
	}
	time.Sleep(time.Millisecond * time.Duration(t))
	return 0
}

func LuaLibPrintf(ls *lua.LState) int {
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
		v[i-2] = LuaToGo(ls.Get(i))
	}
	str := fmt.Sprintf(format, v...)
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		Logger.LogError("common.LuaLibPrintf", "%v", err)
		return 0
	}
	defer fp.Close()
	if _, err := fp.Write([]byte(str)); err != nil {
		Logger.LogError("common.LuaLibPrintf", "%v", err)
		return 0
	}
	return 0
}

func LuaLibNewCron(ls *lua.LState) int {
	expr, err := cronexpr.Parse(ls.ToString(1))
	if err != nil {
		ls.Push(lua.LNil)
		ls.Push(lua.LString(err.Error()))
		return 2
	}
	ud := ls.NewUserData()
	ud.Value = expr
	index := ls.NewTable()
	index.RawSetString("next", ls.NewFunction(luaLibCronNext))
	meta := ls.NewTable()
	meta.RawSetString("__index", index)
	ud.Metatable = meta
	ls.Push(ud)
	ls.Push(lua.LNil)
	return 2
}

func luaLibCronNext(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	ls.Push(lua.LNumber(float64(ud.Value.(*cronexpr.Expression).Next(time.Now()).Unix())))
	return 1
}
