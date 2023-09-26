package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/dbutil"
)

type DeleteTaskController struct {
	BaseController
	Path string
}

func (c *DeleteTaskController) POST(ctx *gin.Context) {
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
		"select name from task where id = ?",
		id,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	if len(rows) == 0 {
		c.Failure(ctx, ErrTaskNotExist)
		return
	}

	name := rows[0]["name"]

	_, err = db.Exec(
		"delete from task where id = ?",
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	common.Logger.LogInfo("apisvr.DeleteTaskController.POST", "task `%s` deleted by `%s`", name, user)

	c.Success(ctx, nil)
}
