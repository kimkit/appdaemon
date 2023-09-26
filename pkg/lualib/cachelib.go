package lualib

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
	"github.com/yuin/gopher-lua"
)

type CacheLib struct {
	clients map[string]gcache.Cache
}

func NewCacheLib(clients map[string]gcache.Cache) *CacheLib {
	if clients == nil {
		clients = make(map[string]gcache.Cache)
	}
	return &CacheLib{clients}
}

func (cl *CacheLib) RegisterGlobal(ls *lua.LState, name string) {
	t := ls.NewTable()
	t.RawSetString("set", ls.NewFunction(cl.Set))
	t.RawSetString("get", ls.NewFunction(cl.Get))
	t.RawSetString("has", ls.NewFunction(cl.Has))
	t.RawSetString("remove", ls.NewFunction(cl.Remove))
	ls.SetGlobal(name, t)
}

func (cl *CacheLib) Set(ls *lua.LState) int {
	name := ls.ToString(1)
	if client, ok := cl.clients[name]; ok {
		key := ls.ToString(2)
		val := LuaToGo(ls.Get(3))
		expiration := ls.ToInt(4)
		var err error
		if expiration > 0 {
			err = client.SetWithExpire(key, val, time.Duration(expiration)*time.Second)
		} else {
			err = client.Set(key, val)
		}
		if err != nil {
			ls.Push(lua.LString(err.Error()))
		} else {
			ls.Push(lua.LNil)
		}
		return 1
	}
	ls.Push(lua.LString(fmt.Sprintf("cache client not exist (%s)", name)))
	return 1
}

func (cl *CacheLib) Get(ls *lua.LState) int {
	name := ls.ToString(1)
	if client, ok := cl.clients[name]; ok {
		key := ls.ToString(2)
		val, err := client.Get(key)
		if err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		ls.Push(GoToLua(ls, val))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("cache client not exist (%s)", name)))
	return 2
}

func (cl *CacheLib) Has(ls *lua.LState) int {
	name := ls.ToString(1)
	if client, ok := cl.clients[name]; ok {
		key := ls.ToString(2)
		ls.Push(lua.LBool(client.Has(key)))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LBool(false))
	ls.Push(lua.LString(fmt.Sprintf("cache client not exist (%s)", name)))
	return 2
}

func (cl *CacheLib) Remove(ls *lua.LState) int {
	name := ls.ToString(1)
	if client, ok := cl.clients[name]; ok {
		key := ls.ToString(2)
		ls.Push(lua.LBool(client.Remove(key)))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LBool(false))
	ls.Push(lua.LString(fmt.Sprintf("cache client not exist (%s)", name)))
	return 2
}
