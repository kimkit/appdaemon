package lualib

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/yuin/gopher-lua"
)

type RedisLib struct {
	clients map[string]*redis.Client
}

func NewRedisLib(clients map[string]*redis.Client) *RedisLib {
	if clients == nil {
		clients = make(map[string]*redis.Client)
	}
	return &RedisLib{clients}
}

func (rl *RedisLib) RegisterGlobal(ls *lua.LState, name string) {
	t := ls.NewTable()
	t.RawSetString("call", ls.NewFunction(rl.Call))
	t.RawSetString("pipeline", ls.NewFunction(rl.Pipeline))
	ls.SetGlobal(name, t)
}

func (rl *RedisLib) Call(ls *lua.LState) int {
	name := ls.ToString(1)
	var args []interface{}
	for i := 2; i <= ls.GetTop(); i++ {
		args = append(args, LuaToGo(ls.Get(i)))
	}
	if client, ok := rl.clients[name]; ok {
		cmd := client.Do(args...)
		if err := cmd.Err(); err != nil {
			if err == redis.Nil {
				ls.Push(lua.LNil)
				ls.Push(lua.LNil)
				return 2
			} else {
				ls.Push(lua.LNil)
				ls.Push(lua.LString(err.Error()))
				return 2
			}
		}
		ls.Push(GoToLua(ls, cmd.Val()))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("redis client not exist (%s)", name)))
	return 2
}

func (rl *RedisLib) Pipeline(ls *lua.LState) int {
	name := ls.ToString(1)
	var cmds [][]interface{}
	_cmds := LuaToGo(ls.ToTable(2))
	if __cmds, ok := _cmds.([]interface{}); ok {
		for _, _cmd := range __cmds {
			if __cmd, ok := _cmd.([]interface{}); ok {
				cmds = append(cmds, __cmd)
			}
		}
	}
	if len(cmds) == 0 {
		ls.Push(lua.LNil)
		ls.Push(lua.LNil)
		return 2
	}
	if client, ok := rl.clients[name]; ok {
		pl := client.Pipeline()
		var res []*redis.Cmd
		for _, cmd := range cmds {
			res = append(res, pl.Do(cmd...))
		}
		if _, err := pl.Exec(); err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		t := ls.NewTable()
		for _, cmd := range res {
			row := ls.NewTable()
			row.RawSetString("val", GoToLua(ls, cmd.Val()))
			if err := cmd.Err(); err != nil {
				if err == redis.Nil {
					row.RawSetString("err", lua.LNil)
				} else {
					row.RawSetString("err", lua.LString(err.Error()))
				}
			} else {
				row.RawSetString("err", lua.LNil)
			}
			t.Append(row)
		}
		ls.Push(t)
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("redis client not exist (%s)", name)))
	return 2
}
