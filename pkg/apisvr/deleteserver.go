package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
)

type DeleteServerController struct {
	BaseController
	Path string
}

func (c *DeleteServerController) POST(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
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

	_, err = db.Exec(
		"delete from server where id = ?",
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	c.Success(ctx, nil)
}
