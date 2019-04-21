package cmdsvr

import (
	"fmt"
	"strings"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/cmdlib"
	"github.com/kimkit/daemon"
	"github.com/kimkit/redsvr"
)

func authHandler(cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if len(common.Config.Passwords) > 0 {
		if err := cmdlib.CheckAuth(conn); err != nil {
			return err
		}
		redsvr.Log("INFO", "(%s) %s %s", cmdlib.GetAuthUser(conn), cmd.Name, strings.Join(args, " "))
	}
	return nil
}

func init() {
	common.Cmdsvr.Register(cmdlib.NewEchoCommand("echo", nil))
	common.Cmdsvr.Register(cmdlib.NewAuthCommand("auth", common.Config.Passwords))
	common.Cmdsvr.Register(cmdlib.NewJobListCommand("job.list", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobStartCommand("job.start", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobStopCommand("job.stop", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobStopAllCommand("job.stopAll", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobCleanCommand("job.clean", authHandler, common.JobManager))
	common.Cmdsvr.Register(newTaskListCommand("task.list", authHandler))
	common.Cmdsvr.Register(newTaskAddCommand("task.add", authHandler))
	common.Cmdsvr.Register(newTaskDeleteCommand("task.delete", authHandler))
}

func Run() {
	if common.Config.Daemon {
		logFile := fmt.Sprintf("%s/appdaemon.log", common.Config.LogsDir)
		pidFile := fmt.Sprintf("%s/appdaemon.pid", common.Config.LogsDir)
		daemon.Daemon(logFile, pidFile)
	}
	common.JobManager.LoadJobs(common.Config.JobsFile, func(info []interface{}) error {
		return common.Client.Do(info...).Err()
	})
	common.Cmdsvr.ListenAndServe(common.Config.Addr)
	common.JobManager.SaveRunningJobs(common.Config.JobsFile)
	common.JobManager.StopAllJobs()
}
