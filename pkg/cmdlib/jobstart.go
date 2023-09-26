package cmdlib

import (
	"fmt"

	"github.com/kimkit/appdaemon/pkg/jobext"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type jobStartCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	Jm          *jobext.JobManager
}

func (cmd *jobStartCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	job := cmd.Jm.GetJob(args[0], nil)
	if job == nil {
		return fmt.Errorf("job `%s` not exist", args[0])
	}
	if err := job.Start(); err != nil {
		return nil
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

func NewJobStartCommand(name string, handler redsvr.CommandHandler, jm *jobext.JobManager) *redsvr.Command {
	return redsvr.NewCommand(&jobStartCommand{Name: name, Argc: 1, AuthHandler: handler, Jm: jm})
}
