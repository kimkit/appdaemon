package cmdlib

import (
	"github.com/kimkit/appdaemon/pkg/jobctl"
	"github.com/kimkit/appdaemon/pkg/jobext"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type jobCleanCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	Jm          *jobext.JobManager
}

func (cmd *jobCleanCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	cmd.Jm.Map.Range(func(k, v interface{}) bool {
		job := v.(*jobctl.Job)
		if !job.IsRunning() {
			cmd.Jm.Map.Delete(k)
		}
		return true
	})
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

func NewJobCleanCommand(name string, handler redsvr.CommandHandler, jm *jobext.JobManager) *redsvr.Command {
	return redsvr.NewCommand(&jobCleanCommand{Name: name, Argc: 0, AuthHandler: handler, Jm: jm})
}
