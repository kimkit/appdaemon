package lualib

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/yuin/gopher-lua"
)

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
	case float32:
		return lua.LNumber(float64(v))
	case int:
		return lua.LNumber(float64(v))
	case int8:
		return lua.LNumber(float64(v))
	case int16:
		return lua.LNumber(float64(v))
	case int32:
		return lua.LNumber(float64(v))
	case int64:
		return lua.LNumber(float64(v))
	case uint:
		return lua.LNumber(float64(v))
	case uint8:
		return lua.LNumber(float64(v))
	case uint16:
		return lua.LNumber(float64(v))
	case uint32:
		return lua.LNumber(float64(v))
	case uint64:
		return lua.LNumber(float64(v))
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

func Sleep(ls *lua.LState) int {
	t := ls.ToInt(1)
	if t < 0 {
		t = 0
	}
	time.Sleep(time.Millisecond * time.Duration(t))
	return 0
}

func NewCron(ls *lua.LState) int {
	expr, err := cronexpr.Parse(ls.ToString(1))
	if err != nil {
		ls.Push(lua.LNil)
		ls.Push(lua.LString(err.Error()))
		return 2
	}
	ud := ls.NewUserData()
	ud.Value = expr
	index := ls.NewTable()
	index.RawSetString("next", ls.NewFunction(cronNext))
	meta := ls.NewTable()
	meta.RawSetString("__index", index)
	ud.Metatable = meta
	ls.Push(ud)
	ls.Push(lua.LNil)
	return 2
}

func cronNext(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	ls.Push(lua.LNumber(float64(ud.Value.(*cronexpr.Expression).Next(time.Now()).Unix())))
	return 1
}

func UUID(ls *lua.LState) int {
	u, err := uuid.NewRandom()
	if err != nil {
		ls.Push(lua.LNil)
		ls.Push(lua.LString(err.Error()))
		return 2
	}
	ls.Push(lua.LString(u.String()))
	ls.Push(lua.LNil)
	return 2
}

func MD5(ls *lua.LState) int {
	h := md5.New()
	io.WriteString(h, ls.ToString(1))
	ls.Push(lua.LString(fmt.Sprintf("%x", h.Sum(nil))))
	return 1
}

func Trim(ls *lua.LState) int {
	str := ls.ToString(1)
	chars := ls.ToString(2)
	flag := ls.ToString(3)
	switch strings.ToLower(flag) {
	case "l", "left":
		if chars == "" {
			ls.Push(lua.LString(strings.TrimLeftFunc(str, unicode.IsSpace)))
		} else {
			ls.Push(lua.LString(strings.TrimLeft(str, chars)))
		}
	case "r", "right":
		if chars == "" {
			ls.Push(lua.LString(strings.TrimRightFunc(str, unicode.IsSpace)))
		} else {
			ls.Push(lua.LString(strings.TrimRight(str, chars)))
		}
	default:
		if chars == "" {
			ls.Push(lua.LString(strings.TrimSpace(str)))
		} else {
			ls.Push(lua.LString(strings.Trim(str, chars)))
		}
	}
	return 1
}

func Split(ls *lua.LState) int {
	arr := strings.Split(ls.ToString(1), ls.ToString(2))
	t := ls.NewTable()
	for _, v := range arr {
		t.Append(lua.LString(v))
	}
	ls.Push(t)
	return 1
}

func Random(ls *lua.LState) int {
	max := ls.ToInt(1)
	val := -1
	if max > 0 {
		val = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
	}
	ls.Push(lua.LNumber(val))
	return 1
}
