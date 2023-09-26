package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/dbutil"
)

type UpdateServerStatusController struct {
	BaseController
	Path string
}

func (c *UpdateServerStatusController) POST(ctx *gin.Context) {
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
		"update server set status = ? where id = ?",
		status,
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	common.Logger.LogInfo("apisvr.UpdateServerStatusController.POST", "server `%s` status `%d` updated by `%s`", addr, status, user)

	c.Success(ctx, nil)
}
