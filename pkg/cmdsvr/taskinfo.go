package cmdsvr

import (
	"strconv"
	"strings"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/jobctl"
	"github.com/kimkit/redsvr"
)

type taskInfoCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *taskInfoCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	key := GetTaskKey(args[0])
	list := []interface{}{}
	common.JobManager.Map.Range(func(k, v interface{}) bool {
		name, _ := k.(string)
		job, _ := v.(*jobctl.Job)
		if job != nil {
			if name == key {
				list = append(list, []interface{}{
					common.Time2str(common.GetMapValueInt(&job.Map, "last")),
					common.GetMapValueString(&job.Map, "name"),
					common.GetMapValueString(&job.Map, "rule"),
					common.GetMapValueStringArr(&job.Map, "args"),
				})
				return true
			}
			if strings.HasPrefix(name, key+"_") {
				if _, err := strconv.Atoi(name[len(key)+1:]); err != nil {
					// ignore
				} else {
					list = append(list, []interface{}{
						common.Time2str(common.GetMapValueInt(&job.Map, "last")),
						common.GetMapValueString(&job.Map, "name"),
						common.GetMapValueString(&job.Map, "rule"),
						common.GetMapValueStringArr(&job.Map, "args"),
					})
					return true
				}
			}
		}
		return true
	})

	redsvr.WriteArray(conn, list)
	return nil
}

func newTaskInfoCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&taskInfoCommand{Name: name, Argc: 1, AuthHandler: handler})
}
