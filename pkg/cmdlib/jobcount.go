package cmdlib

import (
	"github.com/kimkit/appdaemon/pkg/jobext"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type jobCountCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	Jm          *jobext.JobManager
}

func (cmd *jobCountCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	redsvr.WriteInt(conn, len(cmd.Jm.GetRunningJobs()))
	return nil
}

func NewJobCountCommand(name string, handler redsvr.CommandHandler, jm *jobext.JobManager) *redsvr.Command {
	return redsvr.NewCommand(&jobCountCommand{Name: name, Argc: 0, AuthHandler: handler, Jm: jm})
}
