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

	all := ctx.Query("all")

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	sql := ""
	if all == "1" {
		sql = "select id,addr,status,updatetime from server order by updatetime desc"
	} else {
		sql = "select id,addr,status,updatetime from server where status = 1 order by updatetime desc"
	}
	list, err := dbutil.FetchAll(db.Query(sql))
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	c.Success(ctx, gin.H{
		"list": list,
	})
}
