package apisvr

import (
	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
)

type GetServerListController struct {
	BaseController
	Path string
}

func (c *GetServerListController) GET(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
		c.Failure(ctx, err)
		return
	}

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	list, err := dbutil.FetchAll(db.Query("select addr from server where status = 1 order by updatetime desc"))
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	common.SortMaps(list, "addr")

	c.Success(ctx, gin.H{
		"list": list,
	})
}
