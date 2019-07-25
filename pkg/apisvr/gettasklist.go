package apisvr

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
)

type GetTaskListController struct {
	BaseController
	Path string
}

func (c *GetTaskListController) GET(ctx *gin.Context) {
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
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	addr := strings.TrimSpace(ctx.Query("addr"))

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	list, err := dbutil.FetchAll(db.Query(
		"select id,addr,name,rule,command,status from task where name like ? and addr = if(?='',addr,?) order by id desc limit ?,?",
		"%"+keyword+"%",
		addr,
		addr,
		(page-1)*pagesize,
		pagesize,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select count(*) as total from task where name like ? and addr = if(?='',addr,?)",
		"%"+keyword+"%",
		addr,
		addr,
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
