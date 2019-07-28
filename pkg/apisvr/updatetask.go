package apisvr

import (
	"fmt"
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

type UpdateTaskController struct {
	BaseController
	Path string
}

func (c *UpdateTaskController) POST(ctx *gin.Context) {
	user, err := c.GetLoginUser(ctx)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	idStr := ctx.PostForm("id")
	id, _ := strconv.Atoi(idStr)
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
	ruleType := "single"
	processNum := 0
	if rule != "" {
		if _, err := cronexpr.Parse(rule); err != nil {
			num, err := strconv.Atoi(rule)
			if err != nil || num <= 0 {
				c.Failure(ctx, ErrTaskRuleInvalid)
				return
			}
			ruleType = "multi"
			processNum = num
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
		"select id from task where name = ? and id <> ?",
		name,
		id,
	))
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if len(rows) != 0 {
		c.Failure(ctx, ErrTaskNameDuplicate)
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

	rows, err = dbutil.FetchAll(db.Query(
		"select addr,name,rule,status from task where id = ?",
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
	if rows[0]["status"] != "0" {
		c.Failure(ctx, ErrTaskEnable)
		return
	}

	oldTask := rows[0]

	var checkNames []string
	var checkAddrs []string

	oldRuleType := "single"
	oldProcessNum := 0
	if oldTask["rule"] != "" {
		if _, err := cronexpr.Parse(oldTask["rule"]); err != nil {
			num, err := strconv.Atoi(oldTask["rule"])
			if err != nil || num <= 0 {
				c.Failure(ctx, ErrTaskRuleInvalid)
				return
			}
			oldRuleType = "multi"
			oldProcessNum = num
		}
	}
	if oldRuleType == "single" {
		checkNames = append(checkNames, cmdsvr.GetTaskKey(oldTask["name"]))
	} else {
		for i := 0; i < oldProcessNum; i++ {
			checkNames = append(checkNames, cmdsvr.GetTaskKey(fmt.Sprintf("%s_%03d", oldTask["name"], i)))
		}
	}
	if oldTask["addr"] != "" {
		checkAddrs = append(checkAddrs, oldTask["addr"])
	}
	ret, err := IsRunning(checkNames, checkAddrs)
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if ret {
		c.Failure(ctx, ErrJobIsRunning)
		return
	}

	checkNames = nil
	checkAddrs = nil

	if ruleType == "single" {
		checkNames = append(checkNames, cmdsvr.GetTaskKey(name))
	} else {
		for i := 0; i < processNum; i++ {
			checkNames = append(checkNames, cmdsvr.GetTaskKey(fmt.Sprintf("%s_%03d", name, i)))
		}
	}
	if addr != "" {
		checkAddrs = append(checkAddrs, addr)
	}
	ret, err = IsRunning(checkNames, checkAddrs)
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	if ret {
		c.Failure(ctx, ErrJobIsRunning)
		return
	}

	_, err = db.Exec(
		"update task set addr = ?, name = ?, rule = ?, command = ?, status = ?, updatetime = ?, updateuser = ? where id = ?",
		addr,
		name,
		rule,
		command,
		status,
		time.Now().Format("2006-01-02 15:04:05"),
		user,
		id,
	)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	common.Logger.LogInfo("apisvr.UpdateTaskController.POST", "task `%s` updated by `%s`", name, user)

	c.Success(ctx, nil)
}
