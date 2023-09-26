package cmdsvr

import (
	"fmt"
	"time"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/jobctl"
	"github.com/kimkit/appdaemon/pkg/lualib"
	"github.com/kimkit/appdaemon/pkg/redsvr"
	"github.com/yuin/gopher-lua"
)

type luaScriptRunnerCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func GetLuaScriptPrefix() string {
	return "luascript_"
}

func GetLuaScriptKey(name string) string {
	return fmt.Sprintf("%s%s", GetLuaScriptPrefix(), name)
}

func (cmd *luaScriptRunnerCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	if len(args) < 1 {
		return fmt.Errorf("wrong number of arguments for '%s'", _cmd.Name)
	}
	if len(args) > 1 {
		if !taskNameRegexp.MatchString(args[1]) {
			return fmt.Errorf("script tag `%s` invalid", args[1])
		}
	}
	key := ""
	if len(args) > 1 {
		key = GetLuaScriptKey(args[0] + "_" + args[1])
	} else {
		key = GetLuaScriptKey(args[0])
	}
	job := common.JobManager.GetJob(key, &luaScriptRunnerJob{})
	if job == nil {
		return fmt.Errorf("create job `%s` failed", key)
	}
	job.Map.LoadOrStore("name", args[0])
	if len(args) > 1 {
		job.Map.LoadOrStore("tag", args[1])
		job.Map.LoadOrStore("args", args[2:])
	} else {
		job.Map.LoadOrStore("tag", "")
		job.Map.LoadOrStore("args", []string{})
	}
	if err := job.Start(); err != nil {
		return fmt.Errorf("start job `%s` failed", key)
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

type luaScriptRunnerJob struct {
	name   string
	tag    string
	args   []interface{}
	md5sum string
	ls     *lua.LState
}

func (job *luaScriptRunnerJob) InitHandler(_job *jobctl.Job) {
	job.name = common.GetMapValueString(&_job.Map, "name")
	job.tag = common.GetMapValueString(&_job.Map, "tag")
	args := common.GetMapValueStringArr(&_job.Map, "args")
	for _, v := range args {
		job.args = append(job.args, v)
	}
	job.ls = common.CreateStateHandler()
}

func (job *luaScriptRunnerJob) ExitHandler(_job *jobctl.Job) {
	common.JobManager.DestroyJob(_job)
}

func (job *luaScriptRunnerJob) ExecHandler(_job *jobctl.Job) {
	s := common.LuaScriptStore.Get(job.name)
	if s == nil {
		_job.Stop(0)
		return
	}
	if job.md5sum == "" {
		job.md5sum = s.Md5sum
	} else if job.md5sum != s.Md5sum {
		_job.Stop(0)
		return
	}
	job.ls.SetGlobal("jobname", lua.LString(common.JobManager.GetJobName(_job)))
	job.ls.SetGlobal("scriptname", lua.LString(job.name))
	job.ls.SetGlobal("scriptmd5sum", lua.LString(job.md5sum))
	job.ls.SetGlobal("scripttag", lua.LString(job.tag))
	job.ls.SetGlobal("scriptargs", lualib.GoToLua(job.ls, job.args))
	job.ls.Push(job.ls.NewFunctionFromProto(s.Fp))
	if err := job.ls.PCall(0, 0, nil); err != nil {
		common.Logger.LogError("cmdsvr.luaScriptRunnerJob.ExecHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
		time.Sleep(time.Millisecond * 500)
	}
}

func newLuaScriptRunnerCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&luaScriptRunnerCommand{Name: name, Argc: -1, AuthHandler: handler})
}
