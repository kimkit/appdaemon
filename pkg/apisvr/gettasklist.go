package apisvr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/cmdsvr"
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

	rows, err = dbutil.FetchAll(db.Query(
		"select addr from server where status = 1 order by updatetime desc limit 1",
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	defaultAddr := ""
	if len(rows) > 0 {
		defaultAddr = rows[0]["addr"]
	}

	names := make(map[string][]string)
	for _, row := range list {
		if row["addr"] == "" {
			names[defaultAddr] = append(names[defaultAddr], row["name"])
		} else {
			names[row["addr"]] = append(names[row["addr"]], row["name"])
		}
	}

	var addrsFilter []string
	for addr, _ := range names {
		if addr != "" {
			addrsFilter = append(addrsFilter, fmt.Sprintf("'%s'", common.Addslashes(addr)))
		}
	}
	var addrs []string
	if len(addrsFilter) > 0 {
		rows, err := dbutil.FetchAll(db.Query(
			fmt.Sprintf(
				"select addr from server where status = 1 and unix_timestamp() - unix_timestamp(updatetime) <= 60 and addr in (%s)",
				strings.Join(addrsFilter, ","),
			),
		))
		if err != nil {
			c.Failure(ctx, err)
			return
		}
		for _, row := range rows {
			addrs = append(addrs, row["addr"])
		}
	}

	updateTimeList := make(map[string]string)

	password := ""
	if len(common.Config.Passwords) > 0 {
		password = common.Config.Passwords[0]
	}
	for _, addr := range addrs {
		redis := common.NewRedis(&common.RedisConfig{
			Addr:     addr,
			Password: password,
		})

		var params []interface{}
		params = append(params, "task.updatetime")
		for _, name := range names[addr] {
			params = append(params, name)
		}

		res, err := redis.Do(params...).Result()
		if err != nil {
			common.Logger.LogError("apisvr.GetTaskListController.GET", "%v", err)
			continue
		}
		if _res, ok := res.([]interface{}); ok {
			for k, v := range _res {
				if _v, ok := v.(string); ok {
					if _v != "0000-00-00 00:00:00" {
						if len(names[addr]) > k {
							updateTimeList[names[addr][k]] = _v
						}
					}
				}
			}
		}
	}

	for _, row := range list {
		row["updatetime"] = updateTimeList[row["name"]]
		row["jobname"] = cmdsvr.GetTaskKey(row["name"])
		if row["addr"] == "" {
			row["subaddr"] = defaultAddr
		} else {
			row["subaddr"] = row["addr"]
		}
	}

	c.Success(ctx, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pagesize": pagesize,
	})
}
