package lualib

import (
	"encoding/json"
	"fmt"
	"net/http"
	stdurl "net/url"

	"github.com/kimkit/appdaemon/pkg/reqctl"
	"github.com/yuin/gopher-lua"
)

type HttpLib struct {
	clients map[string]*http.Client
}

func NewHttpLib(clients map[string]*http.Client) *HttpLib {
	if clients == nil {
		clients = make(map[string]*http.Client)
	}
	return &HttpLib{clients}
}

func (hl *HttpLib) RegisterGlobal(ls *lua.LState, name string) {
	t := ls.NewTable()
	t.RawSetString("get", ls.NewFunction(hl.Get))
	t.RawSetString("post", ls.NewFunction(hl.Post))
	ls.SetGlobal(name, t)
}

func (hl *HttpLib) Get(ls *lua.LState) int {
	name := ls.ToString(1)
	url := ls.ToString(2)
	data := LuaToGo(ls.Get(3))
	if _data, ok := data.(map[string]interface{}); ok {
		__data := make(map[string]string)
		for k, v := range _data {
			__data[k] = fmt.Sprint(v)
		}
		data = __data
	}
	header := make(map[string]string)
	_header := LuaToGo(ls.Get(4))
	if __header, ok := _header.(map[string]interface{}); ok {
		for k, v := range __header {
			header[k] = fmt.Sprint(v)
		}
	}
	if client, ok := hl.clients[name]; ok {
		res, err := reqctl.Get(client, url, data, header)
		if err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		header := make(map[string]interface{})
		for k, v := range res.Header {
			var row []interface{}
			for _, vv := range v {
				row = append(row, vv)
			}
			header[k] = row
		}
		ls.Push(GoToLua(ls, map[string]interface{}{
			"statuscode": res.StatusCode,
			"header":     header,
			"body":       string(res.Body),
		}))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("http client not exist (%s)", name)))
	return 2
}

func (hl *HttpLib) Post(ls *lua.LState) int {
	name := ls.ToString(1)
	url := ls.ToString(2)
	data := LuaToGo(ls.Get(3))
	if _data, ok := data.(map[string]interface{}); ok {
		__data := make(map[string]string)
		for k, v := range _data {
			__data[k] = fmt.Sprint(v)
		}
		data = __data
	}
	header := make(map[string]string)
	_header := LuaToGo(ls.Get(4))
	if __header, ok := _header.(map[string]interface{}); ok {
		for k, v := range __header {
			header[k] = fmt.Sprint(v)
		}
	}
	if client, ok := hl.clients[name]; ok {
		res, err := reqctl.Post(client, url, data, header)
		if err != nil {
			ls.Push(lua.LNil)
			ls.Push(lua.LString(err.Error()))
			return 2
		}
		header := make(map[string]interface{})
		for k, v := range res.Header {
			var row []interface{}
			for _, vv := range v {
				row = append(row, vv)
			}
			header[k] = row
		}
		ls.Push(GoToLua(ls, map[string]interface{}{
			"statuscode": res.StatusCode,
			"header":     header,
			"body":       string(res.Body),
		}))
		ls.Push(lua.LNil)
		return 2
	}
	ls.Push(lua.LNil)
	ls.Push(lua.LString(fmt.Sprintf("http client not exist (%s)", name)))
	return 2
}

func QueryBuild(ls *lua.LState) int {
	t := LuaToGo(ls.ToTable(1))
	if m, ok := t.(map[string]interface{}); ok {
		v := stdurl.Values{}
		for key, value := range m {
			if vv, ok := value.([]interface{}); ok {
				for _, vvv := range vv {
					v.Add(fmt.Sprint(key), fmt.Sprint(vvv))
				}
			} else {
				v.Set(fmt.Sprint(key), fmt.Sprint(value))
			}
		}
		ls.Push(lua.LString(v.Encode()))
	} else {
		ls.Push(lua.LString(""))
	}
	return 1
}

func QueryParse(ls *lua.LState) int {
	v, _ := stdurl.ParseQuery(ls.ToString(1))
	res := make(map[string]interface{})
	for key, value := range v {
		var vv []interface{}
		for _, vvv := range value {
			vv = append(vv, vvv)
		}
		res[key] = vv
	}
	ls.Push(GoToLua(ls, res))
	return 1
}

func JsonEncode(ls *lua.LState) int {
	v, err := json.Marshal(LuaToGo(ls.Get(1)))
	if err != nil {
		// pass
	}
	ls.Push(lua.LString(string(v)))
	return 1
}

func JsonDecode(ls *lua.LState) int {
	var v interface{}
	if err := json.Unmarshal([]byte(ls.ToString(1)), &v); err != nil {
		// pass
	}
	ls.Push(GoToLua(ls, v))
	return 1
}
