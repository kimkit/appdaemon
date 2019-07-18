package cmdsvr

import (
	"fmt"
	"time"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/jobctl"
	"github.com/kimkit/redsvr"
)

type luaScriptLoaderCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
	RunnerName  string
}

func (cmd *luaScriptLoaderCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	key := fmt.Sprintf("%s", cmd.Name)
	job := common.JobManager.GetJob(key, &luaScriptLoaderJob{})
	if job == nil {
		return fmt.Errorf("create job `%s` failed", key)
	}
	job.Map.LoadOrStore("runnerName", cmd.RunnerName)
	if err := job.Start(); err != nil {
		return fmt.Errorf("start job `%s` failed", key)
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

type luaScriptLoaderJob struct {
	step       int
	lastTime   int64
	runnerName string
}

func (job *luaScriptLoaderJob) InitHandler(_job *jobctl.Job) {
	job.step = 1
	job.lastTime = 0
	job.runnerName = common.GetMapValueString(&_job.Map, "runnerName")
}

func (job *luaScriptLoaderJob) ExecHandler(_job *jobctl.Job) {
	if job.step == 0 {
		if common.Lister == nil {
			common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "lister should not be nil")
			_job.Stop(0)
			return
		}
		rows, err := common.Lister.FetchRows()
		if err != nil {
			common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "%v", err)
			time.Sleep(time.Millisecond * 500)
			return
		}
		for _, _row := range rows {
			if row, ok := _row.(map[string]string); ok {
				if err := common.LuaScriptStore.Add(row["name"], row["script"]); err != nil {
					common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "%v (%s)", err, row["name"])
				} else {
					// common.Logger.LogDebug("cmdsvr.luaScriptLoaderJob.ExecHandler", "script `%v` loaded", row["name"])
					if err := common.RedisClient.Do(job.runnerName, row["name"], "start").Err(); err != nil {
						common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "%v (%s)", err, row["name"])
					}
				}
			}
		}
		if len(rows) == 0 {
			common.LuaScriptStore.Clean()
			job.step = 1
		}
		return
	}
	if job.step == 1 {
		now := time.Now().Unix()
		if now-job.lastTime >= 2 {
			job.step = 0
			job.lastTime = now
		} else {
			time.Sleep(time.Millisecond * 200)
		}
		return
	}
}

func newLuaScriptLoaderCommand(name, runnerName string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&luaScriptLoaderCommand{Name: name, Argc: 1, AuthHandler: handler, RunnerName: runnerName})
}
