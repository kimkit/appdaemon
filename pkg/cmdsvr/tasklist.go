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
			if len(args) == 0 {
				if strings.HasPrefix(name, GetTaskKeyPrefix()) {
					list = append(list, []string{
						common.Time2str(common.GetMapValueInt(&job.Map, "last")),
						"|",
						common.GetMapValueString(&job.Map, "name"),
						"|",
						common.GetMapValueString(&job.Map, "rule"),
						"|",
						common.Args2str(common.GetMapValueStringArr(&job.Map, "args")),
					})
				}
				return true
			}
			for _, filter := range args {
				if strings.HasPrefix(name, GetTaskKey(filter)+"_") || name == GetTaskKey(filter) {
					list = append(list, []string{
						common.Time2str(common.GetMapValueInt(&job.Map, "last")),
						"|",
						common.GetMapValueString(&job.Map, "name"),
						"|",
						common.GetMapValueString(&job.Map, "rule"),
						"|",
						common.Args2str(common.GetMapValueStringArr(&job.Map, "args")),
					})
				}
			}
		}
		return true
	})
	_list := common.BuildTable(list)
	sort.Strings(_list)
	redsvr.WriteArray(conn, _list)
	return nil
}

func newTaskListCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&taskListCommand{Name: name, Argc: -1, AuthHandler: handler})
}
