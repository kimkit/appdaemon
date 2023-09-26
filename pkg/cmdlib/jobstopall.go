package cmdlib

import (
	"github.com/kimkit/appdaemon/pkg/jobctl"
	"github.com/kimkit/appdaemon/pkg/jobext"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type jobStopAllCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	Jm          *jobext.JobManager
}

func (cmd *jobStopAllCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	cmd.Jm.Map.Range(func(k, v interface{}) bool {
		job := v.(*jobctl.Job)
		job.Stop(0)
		return true
	})
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

func NewJobStopAllCommand(name string, handler redsvr.CommandHandler, jm *jobext.JobManager) *redsvr.Command {
	return redsvr.NewCommand(&jobStopAllCommand{Name: name, Argc: 0, AuthHandler: handler, Jm: jm})
}
