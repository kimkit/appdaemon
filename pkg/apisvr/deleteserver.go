package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/dbutil"
)

type DeleteServerController struct {
	BaseController
	Path string
}

func (c *DeleteServerController) POST(ctx *gin.Context) {
	user, err := c.GetLoginUser(ctx)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	idStr := ctx.PostForm("id")
	id, _ := strconv.Atoi(idStr)

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select addr from server where id = ?",
		id,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	if len(rows) == 0 {
		c.Failure(ctx, ErrServerNotExist)
		return
	}

	addr := rows[0]["addr"]

	_, err = db.Exec(
		"delete from server where id = ?",
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	common.Logger.LogInfo("apisvr.DeleteServerController.POST", "server `%s` deleted by `%s`", addr, user)

	c.Success(ctx, nil)
}
