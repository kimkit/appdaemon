package cmdlib

import (
	"fmt"
	"strconv"

	"github.com/kimkit/appdaemon/pkg/jobext"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type jobStopCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	Jm          *jobext.JobManager
}

func (cmd *jobStopCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("wrong number of arguments for '%s'", cmd.Name)
	}
	job := cmd.Jm.GetJob(args[0], nil)
	if job == nil {
		return fmt.Errorf("job `%s` not exist", args[0])
	}
	second := 0
	if len(args) == 2 {
		second, _ = strconv.Atoi(args[1])
	}
	if err := job.Stop(second); err != nil {
		return err
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

func NewJobStopCommand(name string, handler redsvr.CommandHandler, jm *jobext.JobManager) *redsvr.Command {
	return redsvr.NewCommand(&jobStopCommand{Name: name, Argc: -1, AuthHandler: handler, Jm: jm})
}
