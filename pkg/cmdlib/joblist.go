package cmdlib

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kimkit/appdaemon/pkg/jobctl"
	"github.com/kimkit/appdaemon/pkg/jobext"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type jobListCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	Jm          *jobext.JobManager
}

func (cmd *jobListCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	var list []string
	cmd.Jm.Map.Range(func(k, v interface{}) bool {
		name := k.(string)
		job := v.(*jobctl.Job)
		status := 0
		if job.IsRunning() {
			status = 1
		}
		if len(args) == 0 {
			list = append(list, fmt.Sprintf("%d %s", status, name))
		} else {
			for _, filter := range args {
				if strings.Contains(name, filter) {
					list = append(list, fmt.Sprintf("%d %s", status, name))
					break
				}
			}
		}
		return true
	})
	sort.Strings(list)
	redsvr.WriteArray(conn, list)
	return nil
}

func NewJobListCommand(name string, handler redsvr.CommandHandler, jm *jobext.JobManager) *redsvr.Command {
	return redsvr.NewCommand(&jobListCommand{Name: name, Argc: -1, AuthHandler: handler, Jm: jm})
}
