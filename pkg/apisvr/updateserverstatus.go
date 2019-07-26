package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
)

type UpdateServerStatusController struct {
	BaseController
	Path string
}

func (c *UpdateServerStatusController) POST(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
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

	_, err = db.Exec(
		"update server set status = ? where id = ?",
		status,
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	c.Success(ctx, nil)
}
