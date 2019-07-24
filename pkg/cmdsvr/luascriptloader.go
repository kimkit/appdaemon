package cmdsvr

import (
	"fmt"
	"strings"
	"time"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/jobctl"
	"github.com/kimkit/lister"
	"github.com/kimkit/redsvr"
	"github.com/yuin/gopher-lua"
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
	ls         lister.Lister
	addr       string
}

func (job *luaScriptLoaderJob) InitHandler(_job *jobctl.Job) {
	job.step = 1
	job.lastTime = 0
	job.runnerName = common.GetMapValueString(&_job.Map, "runnerName")
	ls := common.CreateStateHandler()
	lv := ls.GetGlobal("getserveraddr")
	if _, ok := lv.(*lua.LFunction); ok {
		ls.Push(lv)
		if err := ls.PCall(0, 1, nil); err != nil {
			common.Logger.LogError("cmdsvr.luaScriptLoaderJob.InitHandler", "%v", err)
		} else {
			job.addr = ls.ToString(1)
		}
	}
	if job.addr == "" {
		common.Logger.LogWarning("cmdsvr.luaScriptLoaderJob.InitHandler", "cannot get server addr")
	}
}

func (job *luaScriptLoaderJob) ExecHandler(_job *jobctl.Job) {
	if job.step == 0 {
		if common.DBClient == nil {
			common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "db client should not be nil")
			_job.Stop(0)
			return
		}
		if job.ls == nil {
			db, err := common.DBClient.Open()
			if err != nil {
				common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "%v", err)
				time.Sleep(time.Millisecond * 500)
				return
			}
			job.ls = lister.NewDBLister(
				db,
				common.Config.LuaScript.Sql,
				common.Config.LuaScript.IdName,
				common.Config.LuaScript.IdInit,
			)
		}
		rows, err := job.ls.FetchRows()
		if err != nil {
			common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "%v", err)
			time.Sleep(time.Millisecond * 500)
			return
		}
		for _, _row := range rows {
			if row, ok := _row.(map[string]string); ok {
				if err := common.LuaScriptStore.Add(row["name"], "-- "+row["addr"]+"\n"+row["script"]); err != nil {
					common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "%v (%s)", err, row["name"])
				} else {
					if row["addr"] == "" || row["addr"] == job.addr {
						if strings.HasPrefix(row["name"], common.Config.LuaScript.FilterPrefix) {
							// common.Logger.LogDebug("cmdsvr.luaScriptLoaderJob.ExecHandler", "script `%v` loaded", row["name"])
							if err := common.RedisClient.Do(job.runnerName, row["name"]).Err(); err != nil {
								common.Logger.LogError("cmdsvr.luaScriptLoaderJob.ExecHandler", "%v (%s)", err, row["name"])
							}
						}
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
