package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
)

type UpdateLuaScriptStatusController struct {
	BaseController
	Path string
}

func (c *UpdateLuaScriptStatusController) POST(ctx *gin.Context) {
	user, err := c.GetLoginUser(ctx)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	idStr := ctx.PostForm("id")
	id, _ := strconv.Atoi(idStr)
	statusStr := ctx.PostForm("status")
	status, _ := strconv.Atoi(statusStr)
	if status != 1 {
		status = 0
	}

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select name from luascript where id = ?",
		id,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	if len(rows) == 0 {
		c.Failure(ctx, ErrLuaScriptNotExist)
		return
	}

	name := rows[0]["name"]

	_, err = db.Exec(
		"update luascript set status = ? where id = ?",
		status,
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	common.Logger.LogInfo("apisvr.UpdateLuaScriptStatusController.POST", "luascript `%s` status `%d` updated by `%s`", name, status, user)

	c.Success(ctx, nil)
}
