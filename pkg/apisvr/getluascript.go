package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/dbutil"
)

type GetLuaScriptController struct {
	BaseController
	Path string
}

func (c *GetLuaScriptController) GET(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
		c.Failure(ctx, err)
		return
	}

	idStr := ctx.Query("id")
	id, _ := strconv.Atoi(idStr)

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select id,addr,name,script,status from luascript where id = ?",
		id,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if len(rows) != 1 {
		c.Failure(ctx, ErrLuaScriptNotExist)
		return
	}

	c.Success(ctx, gin.H{
		"id":     rows[0]["id"],
		"addr":   rows[0]["addr"],
		"name":   rows[0]["name"],
		"script": rows[0]["script"],
		"status": rows[0]["status"],
	})
}
