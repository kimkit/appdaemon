package lualib

import (
	"database/sql"
	"fmt"

	"github.com/kimkit/appdaemon/pkg/dbutil"
	"github.com/yuin/gopher-lua"
)

type DBLib struct {
	clients map[string]*dbutil.DBWrapper
}

func NewDBLib(clients map[string]*dbutil.DBWrapper) *DBLib {
	if clients == nil {
		clients = make(map[string]*dbutil.DBWrapper)
	}
	return &DBLib{clients}
}

func (dl *DBLib) RegisterGlobal(ls *lua.LState, name string) {
	t := ls.NewTable()
	t.RawSetString("query", ls.NewFunction(dl.Query))
	t.RawSetString("exec", ls.NewFunction(dl.Exec))
	t.RawSetString("insert", ls.NewFunction(dl.Insert))
	t.RawSetString("begin", ls.NewFunction(dl.Begin))
	ls.SetGlobal(name, t)
}

func (dl *DBLib) Query(ls *lua.LState) int {
	name := ls.ToString(1)
	sql := ls.ToString(2)
	var args []interface{}
	for i := 3; i <= ls.GetTop(); i++ {
		args = append(args, LuaToGo(ls.Get(i)))
	}
	if client, ok := dl.clients[name]; ok {
		rows, err := dbutil.FetchAll(client.Query(sql, args...))
		if err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		var _rows []interface{}
		for _, row := range rows {
			_row := make(map[string]interface{})
			for k, v := range row {
				_row[k] = v
			}
			_rows = append(_rows, _row)
		}
		ls.Push(GoToLua(ls, _rows))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("db client not exist (%s)", name)))
	return 2
}

func (dl *DBLib) Exec(ls *lua.LState) int {
	name := ls.ToString(1)
	sql := ls.ToString(2)
	var args []interface{}
	for i := 3; i <= ls.GetTop(); i++ {
		args = append(args, LuaToGo(ls.Get(i)))
	}
	if client, ok := dl.clients[name]; ok {
		ret, err := dbutil.RowsAffected(client.Exec(sql, args...))
		if err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		ls.Push(GoToLua(ls, ret))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("db client not exist (%s)", name)))
	return 2
}

func (dl *DBLib) Insert(ls *lua.LState) int {
	name := ls.ToString(1)
	sql := ls.ToString(2)
	var args []interface{}
	for i := 3; i <= ls.GetTop(); i++ {
		args = append(args, LuaToGo(ls.Get(i)))
	}
	if client, ok := dl.clients[name]; ok {
		ret, err := dbutil.LastInsertId(client.Exec(sql, args...))
		if err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		ls.Push(GoToLua(ls, ret))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("db client not exist (%s)", name)))
	return 2
}

func (dl *DBLib) Begin(ls *lua.LState) int {
	name := ls.ToString(1)
	if client, ok := dl.clients[name]; ok {
		tx, err := client.Begin()
		if err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		ud := ls.NewUserData()
		ud.Value = tx
		index := ls.NewTable()
		index.RawSetString("query", ls.NewFunction(dl.txQuery))
		index.RawSetString("exec", ls.NewFunction(dl.txExec))
		index.RawSetString("insert", ls.NewFunction(dl.txInsert))
		index.RawSetString("commit", ls.NewFunction(dl.txCommit))
		index.RawSetString("rollback", ls.NewFunction(dl.txRollback))
		meta := ls.NewTable()
		meta.RawSetString("__index", index)
		ud.Metatable = meta
		ls.Push(ud)
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("db client not exist (%s)", name)))
	return 2
}

func (dl *DBLib) txQuery(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	tx := ud.Value.(*sql.Tx)
	sql := ls.ToString(2)
	var args []interface{}
	for i := 3; i <= ls.GetTop(); i++ {
		args = append(args, LuaToGo(ls.Get(i)))
	}
	rows, err := dbutil.FetchAll(tx.Query(sql, args...))
	if err != nil {
		ls.Push(lua.LNil)
		ls.Push(lua.LString(err.Error()))
		return 2
	}
	var _rows []interface{}
	for _, row := range rows {
		_row := make(map[string]interface{})
		for k, v := range row {
			_row[k] = v
		}
		_rows = append(_rows, _row)
	}
	ls.Push(GoToLua(ls, _rows))
	ls.Push(lua.LNil)
	return 2
}

func (dl *DBLib) txExec(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	tx := ud.Value.(*sql.Tx)
	sql := ls.ToString(2)
	var args []interface{}
	for i := 3; i <= ls.GetTop(); i++ {
		args = append(args, LuaToGo(ls.Get(i)))
	}
	ret, err := dbutil.RowsAffected(tx.Exec(sql, args...))
	if err != nil {
		ls.Push(lua.LNil)
		ls.Push(lua.LString(err.Error()))
		return 2
	}
	ls.Push(GoToLua(ls, ret))
	ls.Push(lua.LNil)
	return 2
}

func (dl *DBLib) txInsert(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	tx := ud.Value.(*sql.Tx)
	sql := ls.ToString(2)
	var args []interface{}
	for i := 3; i <= ls.GetTop(); i++ {
		args = append(args, LuaToGo(ls.Get(i)))
	}
	ret, err := dbutil.LastInsertId(tx.Exec(sql, args...))
	if err != nil {
		ls.Push(lua.LNil)
		ls.Push(lua.LString(err.Error()))
		return 2
	}
	ls.Push(GoToLua(ls, ret))
	ls.Push(lua.LNil)
	return 2
}

func (dl *DBLib) txCommit(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	tx := ud.Value.(*sql.Tx)
	if err := tx.Commit(); err != nil {
		ls.Push(lua.LString(err.Error()))
		return 1
	}
	ls.Push(lua.LNil)
	return 1
}

func (dl *DBLib) txRollback(ls *lua.LState) int {
	ud := ls.CheckUserData(1)
	tx := ud.Value.(*sql.Tx)
	if err := tx.Rollback(); err != nil {
		ls.Push(lua.LString(err.Error()))
		return 1
	}
	ls.Push(lua.LNil)
	return 1
}
