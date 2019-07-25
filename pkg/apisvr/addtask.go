package apisvr

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"
	"github.com/kimkit/appdaemon/pkg/cmdsvr"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
	"github.com/mattn/go-shellwords"
)

type AddTaskController struct {
	BaseController
	Path string
}

func (c *AddTaskController) POST(ctx *gin.Context) {
	user, err := c.GetLoginUser(ctx)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	addr := strings.TrimSpace(ctx.PostForm("addr"))
	name := strings.TrimSpace(ctx.PostForm("name"))
	rule := strings.TrimSpace(ctx.PostForm("rule"))
	command := strings.TrimSpace(ctx.PostForm("command"))
	statusStr := ctx.PostForm("status")
	status, _ := strconv.Atoi(statusStr)
	if status != 1 {
		status = 0
	}

	if name == "" {
		c.Failure(ctx, ErrTaskNameEmpty)
		return
	}
	if !cmdsvr.CheckTaskName(name) {
		c.Failure(ctx, ErrTaskNameFormatInvalid)
		return
	}
	if rule == "" {
		c.Failure(ctx, ErrTaskRuleEmpty)
		return
	}
	if _, err := cronexpr.Parse(rule); err != nil {
		num, err := strconv.Atoi(rule)
		if err != nil || num <= 0 {
			c.Failure(ctx, ErrTaskRuleInvalid)
			return
		}
	}
	if command == "" {
		c.Failure(ctx, ErrTaskCommandEmpty)
		return
	}
	if _, err := shellwords.Parse(command); err != nil {
		c.Failure(ctx, ErrTaskCommandInvalid)
		return
	}

	db, err := common.GetDB("#")
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select id from task where name = ?",
		name,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if len(rows) != 0 {
		c.Failure(ctx, ErrTaskExist)
		return
	}

	if addr != "" {
		rows, err = dbutil.FetchAll(db.Query(
			"select id from server where addr = ? and status = 1",
			addr,
		))
		if err != nil {
			c.Failure(ctx, err)
			return
		}
		if len(rows) == 0 {
			c.Failure(ctx, ErrServerAddrNotExist)
			return
		}
	}

	_, err = db.Exec(
		"insert into task (addr,name,rule,command,status,createtime,createuser) values (?, ?, ?, ?, ?, ?, ?)",
		addr,
		name,
		rule,
		command,
		status,
		time.Now().Format("2006-01-02 15:04:05"),
		user,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	c.Success(ctx, nil)
}
