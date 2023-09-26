package cmdsvr

import (
	"github.com/kimkit/appdaemon/pkg/cmdlib"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/daemon"
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

func authHandler(cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	user := "/"
	if len(common.Config.Passwords) > 0 {
		if err := cmdlib.CheckAuth(conn); err != nil {
			return err
		}
		user = cmdlib.GetAuthUser(conn)
	}

	common.Logger.LogInfo("cmdsvr.authHandler", "(%s) %s %s", user, cmd.Name, common.Args2str(args))
	return nil
}

func init() {
	common.CmdSvr.Register(cmdlib.NewPingCommand("ping", nil))
	common.CmdSvr.Register(cmdlib.NewEchoCommand("echo", nil))
	common.CmdSvr.Register(cmdlib.NewAuthCommand("auth", common.Config.Passwords))
	common.CmdSvr.Register(cmdlib.NewJobListCommand("job.list", authHandler, common.JobManager))
	common.CmdSvr.Register(cmdlib.NewJobCountCommand("job.count", authHandler, common.JobManager))
	common.CmdSvr.Register(cmdlib.NewJobStartCommand("job.start", authHandler, common.JobManager))
	common.CmdSvr.Register(cmdlib.NewJobStopCommand("job.stop", authHandler, common.JobManager))
	common.CmdSvr.Register(cmdlib.NewJobStopAllCommand("job.stopAll", authHandler, common.JobManager))
	common.CmdSvr.Register(cmdlib.NewJobCleanCommand("job.clean", authHandler, common.JobManager))
	common.CmdSvr.Register(cmdlib.NewJobStatusCommand("job.status", authHandler, common.JobManager))
	common.CmdSvr.Register(newTaskListCommand("task.list", authHandler))
	common.CmdSvr.Register(newTaskAddCommand("task.add", authHandler))
	common.CmdSvr.Register(newTaskDeleteCommand("task.delete", authHandler))
	common.CmdSvr.Register(newTaskUpdateTimeCommand("task.updatetime", authHandler))
	common.CmdSvr.Register(newLuaScriptRunnerCommand("luascript.runner", authHandler))
	common.CmdSvr.Register(newLuaScriptLoaderCommand("luascript.loader", "luascript.runner", authHandler))
	common.CmdSvr.Register(newSubscribeCommand("subscribe", authHandler))
}

func Run() {
	if common.Config.Daemon {
		daemon.Daemon(common.Config.LogFile, common.Config.PidFile)
	}
	jobInfos := common.GetTaskInfos()
	if common.Config.LuaScript.Enable {
		jobInfos = append(jobInfos, []interface{}{"luascript.loader", "start"})
	}
	common.JobManager.LoadJobs(common.Config.JobsFile, func(info []interface{}) error {
		return common.RedisClient.Do(info...).Err()
	}, jobInfos...)
	common.CmdSvr.ListenAndServe(common.Config.Addr)
	common.Logger.LogInfo("cmdsvr.Run", "save running jobs ...")
	common.JobManager.SaveRunningJobs(common.Config.JobsFile)
	common.Logger.LogInfo("cmdsvr.Run", "stop running jobs ...")
	common.JobManager.StopAllJobs()
	common.Logger.LogInfo("cmdsvr.Run", "exit")
}
