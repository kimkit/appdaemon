package cmdsvr

import (
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/redsvr"
)

type taskDeleteCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *taskDeleteCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	key := getTaskKey(args[0])
	job := common.JobManager.GetJob(key, nil)
	if job != nil {
		if err := job.Stop(0); err != nil {
			return err
		}
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

func newTaskDeleteCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&taskDeleteCommand{Name: name, Argc: -1, AuthHandler: handler})
}
