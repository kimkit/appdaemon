package cmdlib

import (
	"fmt"

	"github.com/kimkit/appdaemon/pkg/jobext"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type jobStatusCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	Jm          *jobext.JobManager
}

func (cmd *jobStatusCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	if len(args) < 1 {
		return fmt.Errorf("wrong number of arguments for '%s'", cmd.Name)
	}

	var list []string
	for _, name := range args {
		job := cmd.Jm.GetJob(name, nil)
		if job != nil {
			status := 0
			if job.IsRunning() {
				status = 1
			}
			list = append(list, fmt.Sprintf("%d %s", status, name))
		} else {
			list = append(list, fmt.Sprintf("%d %s", 2, name))
		}
	}

	redsvr.WriteArray(conn, list)
	return nil
}

func NewJobStatusCommand(name string, handler redsvr.CommandHandler, jm *jobext.JobManager) *redsvr.Command {
	return redsvr.NewCommand(&jobStatusCommand{Name: name, Argc: -1, AuthHandler: handler, Jm: jm})
}
