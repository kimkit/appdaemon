package apisvr

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/apires"
	"github.com/kimkit/appdaemon/pkg/cmdsvr"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
	"github.com/kimkit/luactl"
)

type AddLuaScriptController struct {
	BaseController
	Path string
}

func (c *AddLuaScriptController) POST(ctx *gin.Context) {
	user, err := c.GetLoginUser(ctx)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	addr := strings.TrimSpace(ctx.PostForm("addr"))
	name := strings.TrimSpace(ctx.PostForm("name"))
	script := strings.TrimSpace(ctx.PostForm("script"))
	statusStr := ctx.PostForm("status")
	status, _ := strconv.Atoi(statusStr)
	if status != 1 {
		status = 0
	}

	if name == "" {
		c.Failure(ctx, ErrLuaScriptNameEmpty)
		return
	}
	if err := luactl.CheckLuaScriptName(name); err != nil {
		c.Failure(ctx, ErrLuaScriptNameFormatInvalid)
		return
	}
	if script == "" {
		c.Failure(ctx, ErrLuaScriptEmpty)
		return
	}
	if _, err := luactl.DefaultCompileHandler(name, script); err != nil {
		_err := ErrLuaScriptSyntaxError.(*apires.Error).Clone()
		_err.(*apires.Error).Reply.Data = gin.H{
			"error": err.Error(),
		}
		c.Failure(ctx, _err)
		return
	}

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select id from luascript where name = ?",
		name,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if len(rows) != 0 {
		c.Failure(ctx, ErrLuaScriptExist)
		return
	}

	if addr != "" {
		rows, err = dbutil.FetchAll(db.Query(
			"select id from server where addr = ? and status = 1",
			addr,
		))
		if err != nil {
			c.Failure(ctx, err)
			return
		}
		if len(rows) == 0 {
			c.Failure(ctx, ErrServerAddrNotExist)
			return
		}
	}

	var checkNames []string
	var checkAddrs []string
	checkNames = append(checkNames, cmdsvr.GetLuaScriptKey(name))
	if addr != "" {
		checkAddrs = append(checkAddrs, addr)
	}
	ret, err := IsRunning(checkNames, checkAddrs)
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if ret {
		c.Failure(ctx, ErrJobIsRunning)
		return
	}

	_, err = db.Exec(
		"insert into luascript (addr,name,script,status,createtime,createuser) values (?, ?, ?, ?, ?, ?)",
		addr,
		name,
		script,
		status,
		time.Now().Format("2006-01-02 15:04:05"),
		user,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	common.Logger.LogInfo("apisvr.AddLuaScriptController.POST", "luascript `%s` added by `%s`", name, user)

	c.Success(ctx, nil)
}
