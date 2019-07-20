package apisvr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
)

type GetLuaScriptListController struct {
	BaseController
	Path string
}

func (c *GetLuaScriptListController) GET(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
		c.Failure(ctx, err)
		return
	}

	pageStr := ctx.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pagesizeStr := ctx.Query("pagesize")
	if pagesizeStr == "" {
		pagesizeStr = "20"
	}
	pagesize, _ := strconv.Atoi(pagesizeStr)
	if pagesize < 1 {
		pagesize = 1
	}
	keyword := ctx.Query("keyword")

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	list, err := dbutil.FetchAll(db.Query(
		"select id,name,script,status from luascript where name like ? order by id desc limit ?,?",
		"%"+keyword+"%",
		(page-1)*pagesize,
		pagesize,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select count(*) as total from luascript where name like ?",
		"%"+keyword+"%",
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	total := 0
	if len(rows) == 1 {
		total, _ = strconv.Atoi(rows[0]["total"])
	}

	c.Success(ctx, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pagesize": pagesize,
	})
}
