package cmdsvr

import (
	"sort"
	"strings"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/jobctl"
	"github.com/kimkit/redsvr"
)

type taskListCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *taskListCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	var list [][]string
	common.JobManager.Map.Range(func(k, v interface{}) bool {
		name, _ := k.(string)
		job, _ := v.(*jobctl.Job)
		if job != nil {
			if strings.HasPrefix(name, GetTaskKeyPrefix()) {
				list = append(list, []string{
					common.Time2str(common.GetMapValueInt(&job.Map, "last")),
					"|",
					common.GetMapValueString(&job.Map, "name"),
					"|",
					common.GetMapValueString(&job.Map, "rule"),
					"|",
					strings.Join(common.GetMapValueStringArr(&job.Map, "args"), " "),
				})
			}
		}
		return true
	})
	_list := common.BuildTable(list)
	for _, filter := range args {
		var __list []string
		for _, line := range _list {
			if strings.Contains(line, filter) {
				__list = append(__list, line)
			}
		}
		_list = __list
	}
	sort.Strings(_list)
	redsvr.WriteArray(conn, _list)
	return nil
}

func newTaskListCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&taskListCommand{Name: name, Argc: -1, AuthHandler: handler})
}
