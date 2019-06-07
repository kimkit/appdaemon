package cmdsvr

import (
	"fmt"
	"time"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/jobctl"
	"github.com/kimkit/redsvr"
)

type luaScriptRunnerCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func getLuaScriptPrefix() string {
	return "luascript_"
}

func getLuaScriptKey(name string) string {
	return fmt.Sprintf("%s%s", getLuaScriptPrefix(), name)
}

func (cmd *luaScriptRunnerCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	key := getLuaScriptKey(args[0])
	job := common.JobManager.GetJob(key, &luaScriptRunnerJob{})
	if job == nil {
		return fmt.Errorf("create job `%s` failed", key)
	}
	job.Map.LoadOrStore("name", args[0])
	if err := job.Start(); err != nil {
		return fmt.Errorf("start job `%s` failed", key)
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

type luaScriptRunnerJob struct {
	name string
}

func (job *luaScriptRunnerJob) InitHandler(_job *jobctl.Job) {
	job.name = common.GetMapValueString(&_job.Map, "name")
}

func (job *luaScriptRunnerJob) ExecHandler(_job *jobctl.Job) {
	if err := common.LuaScriptStore.Run(job.name); err != nil {
		common.Logger.LogError("cmdsvr.luaScriptRunnerJob.ExecHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
		time.Sleep(time.Millisecond * 500)
	}
}

func newLuaScriptRunnerCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&luaScriptRunnerCommand{Name: name, Argc: 2, AuthHandler: handler})
}
