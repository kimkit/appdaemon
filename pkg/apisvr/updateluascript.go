package apisvr

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/apires"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
	"github.com/kimkit/luactl"
)

type UpdateLuaScriptController struct {
	BaseController
	Path string
}

func (c *UpdateLuaScriptController) POST(ctx *gin.Context) {
	user, err := c.GetLoginUser(ctx)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	idStr := ctx.PostForm("id")
	id, _ := strconv.Atoi(idStr)
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
		"select id from luascript where name = ? and id <> ?",
		name,
		id,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if len(rows) != 0 {
		c.Failure(ctx, ErrLuaScriptNameDuplicate)
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

	_, err = db.Exec(
		"update luascript set addr = ?, name = ?, script = ?, status = ?, updatetime = ?, updateuser = ? where id = ?",
		addr,
		name,
		script,
		status,
		time.Now().Format("2006-01-02 15:04:05"),
		user,
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	c.Success(ctx, nil)
}
