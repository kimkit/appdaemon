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
		if cmd.Name == "luascript.runner" {
			return nil
		}
		var _args []string
		for _, arg := range args {
			if arg == "" {
				arg = "\"\""
			} else if strings.Contains(arg, " ") {
				arg = fmt.Sprintf("\"%s\"", arg)
			}
			_args = append(_args, arg)
		}
		common.Logger.LogInfo("cmdsvr.authHandler", "(%s) %s %s", cmdlib.GetAuthUser(conn), cmd.Name, strings.Join(_args, " "))
	}
	return nil
}

func init() {
	common.Cmdsvr.Register(cmdlib.NewPingCommand("ping", nil))
	common.Cmdsvr.Register(cmdlib.NewEchoCommand("echo", nil))
	common.Cmdsvr.Register(cmdlib.NewAuthCommand("auth", common.Config.Passwords))
	common.Cmdsvr.Register(cmdlib.NewJobListCommand("job.list", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobCountCommand("job.count", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobStartCommand("job.start", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobStopCommand("job.stop", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobStopAllCommand("job.stopAll", authHandler, common.JobManager))
	common.Cmdsvr.Register(cmdlib.NewJobCleanCommand("job.clean", authHandler, common.JobManager))
	common.Cmdsvr.Register(newTaskListCommand("task.list", authHandler))
	common.Cmdsvr.Register(newTaskAddCommand("task.add", authHandler))
	common.Cmdsvr.Register(newTaskDeleteCommand("task.delete", authHandler))
	common.Cmdsvr.Register(newLuaScriptRunnerCommand("luascript.runner", authHandler))
	common.Cmdsvr.Register(newLuaScriptLoaderCommand("luascript.loader", "luascript.runner", authHandler))
}

func Run() {
	if common.Config.Daemon {
		daemon.Daemon(common.Config.LogFile, common.Config.PidFile)
	}
	jobInfos := common.GetTaskInfos()
	if common.Lister != nil {
		jobInfos = append(jobInfos, []interface{}{"luascript.loader", "start"})
	}
	common.JobManager.LoadJobs(common.Config.JobsFile, func(info []interface{}) error {
		return common.RedisClient.Do(info...).Err()
	}, jobInfos...)
	common.Cmdsvr.ListenAndServe(common.Config.Addr)
	common.Logger.LogInfo("cmdsvr.Run", "save running jobs ...")
	common.JobManager.SaveRunningJobs(common.Config.JobsFile)
	common.Logger.LogInfo("cmdsvr.Run", "stop running jobs ...")
	common.JobManager.StopAllJobs()
	common.Logger.LogInfo("cmdsvr.Run", "exit")
}
