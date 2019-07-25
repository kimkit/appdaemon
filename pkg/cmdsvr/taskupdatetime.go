package cmdsvr

import (
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/redsvr"
)

type taskUpdateTimeCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *taskUpdateTimeCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	var list []string
	for _, name := range args {
		key := GetTaskKey(name)
		job := common.JobManager.GetJob(key, nil)
		if job != nil {
			list = append(list, common.Time2str(common.GetMapValueInt(&job.Map, "last")))
			continue
		}
		job = common.JobManager.GetJob(key+"_000", nil)
		if job != nil {
			list = append(list, common.Time2str(common.GetMapValueInt(&job.Map, "last")))
			continue
		}
		list = append(list, "0000-00-00 00:00:00")
	}

	redsvr.WriteArray(conn, list)
	return nil
}

func newTaskUpdateTimeCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&taskUpdateTimeCommand{Name: name, Argc: -1, AuthHandler: handler})
}
