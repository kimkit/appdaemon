package cmdsvr

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/jobctl"
	"github.com/kimkit/proctl"
	"github.com/kimkit/redsvr"
)

var (
	taskNameRegexp   = regexp.MustCompile(`^[\w\.\-:]+$`)
	taskSuffixRegexp = regexp.MustCompile(`_\d{3,}$`)
)

type taskAddCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func getTaskKeyPrefix() string {
	return "task_"
}

func getTaskKey(name string) string {
	return fmt.Sprintf("%s%s", getTaskKeyPrefix(), name)
}

func (cmd *taskAddCommand) getJobInfo(args []string) []interface{} {
	_args := []interface{}{cmd.Name}
	for _, arg := range args {
		_args = append(_args, arg)
	}
	return _args
}

func (cmd *taskAddCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	if len(args) < 3 {
		return fmt.Errorf("wrong number of arguments for '%s'", _cmd.Name)
	}
	taskName := args[0]
	taskRule := args[1]
	taskArgs := args[2:]
	if !taskNameRegexp.MatchString(taskName) {
		return fmt.Errorf("task name `%s` invalid", taskName)
	}
	taskExpr, err := cronexpr.Parse(taskRule)
	if err != nil && taskRule != "" {
		num, err := strconv.Atoi(taskRule)
		if err != nil || num < 1 {
			return fmt.Errorf("task rule `%s` invalid", taskRule)
		} else {
			var errs []string
			for i := 0; i < num; i++ {
				var params []interface{}
				params = append(params, cmd.Name)
				params = append(params, fmt.Sprintf("%s_%03d", taskName, i))
				params = append(params, "")
				for j := 2; j < len(args); j++ {
					params = append(params, args[j])
				}
				if err := common.Client.Do(params...).Err(); err != nil {
					if err.Error() != fmt.Sprintf("job `%s_%03d` exist", taskName, i) {
						errs = append(errs, err.Error())
					}
				}
			}
			if len(errs) > 0 {
				return fmt.Errorf("%s", strings.Join(errs, "\n"))
			}
			redsvr.WriteSimpleString(conn, "OK")
			return nil
		}
	}
	key := getTaskKey(taskName)
	job := common.JobManager.GetJob(key, nil)
	if job != nil {
		return fmt.Errorf("job `%s` exist", key)
	}
	job = common.JobManager.GetJob(key, &taskJob{})
	if job == nil {
		return fmt.Errorf("create job `%s` failed", key)
	}
	job.Map.LoadOrStore("name", taskName)
	job.Map.LoadOrStore("rule", taskRule)
	job.Map.LoadOrStore("args", taskArgs)
	job.Map.LoadOrStore("expr", taskExpr)
	job.Map.LoadOrStore("last", 0)

	got := false
	for _, task := range common.Config.Tasks {
		if task.Name == taskName {
			got = true
			break
		}
		if taskSuffixRegexp.MatchString(taskName) && task.Name == taskName[0:len(taskName)-4] {
			got = true
			break
		}
	}
	if !got {
		common.JobManager.SetJobInfo(job, cmd.getJobInfo(args)...)
	}

	if err := job.Start(); err != nil {
		return fmt.Errorf("start job `%s` failed", key)
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

type taskJob struct {
	name string
	rule string
	args []string
	expr *cronexpr.Expression
	next time.Time
	proc *proctl.Process
}

func getTaskOutputFile(job *jobctl.Job) string {
	return fmt.Sprintf(
		"%s%c%s.output",
		common.Config.LogsDir,
		os.PathSeparator,
		common.JobManager.GetJobName(job),
	)
}

func getTaskStatusFile(job *jobctl.Job) string {
	return fmt.Sprintf(
		"%s%c%s.status",
		common.Config.LogsDir,
		os.PathSeparator,
		common.JobManager.GetJobName(job),
	)
}

func (job *taskJob) InitHandler(_job *jobctl.Job) {
	job.name = common.GetMapValueString(&_job.Map, "name")
	job.rule = common.GetMapValueString(&_job.Map, "rule")
	job.args = common.GetMapValueStringArr(&_job.Map, "args")
	expr, _ := _job.Map.Load("expr")
	job.expr, _ = expr.(*cronexpr.Expression)
	if job.expr != nil {
		job.next = job.expr.Next(time.Now())
	}
	job.proc = proctl.NewProcess(job.args[0], job.args[1:]...)
	job.proc.SetOutput(getTaskOutputFile(_job))
	job.proc.SetEnv("TASK_NAME", job.name)
	job.proc.SetEnv("TASK_RULE", job.rule)
	job.proc.SetEnv("TASK_COMMAND", strings.Join(job.args, " "))
	job.proc.SetEnv("TASK_STATUS_FILE", getTaskStatusFile(_job))
}

func (job *taskJob) ExecHandler(_job *jobctl.Job) {
	defer time.Sleep(time.Millisecond * 100)
	now := time.Now()
	if job.expr != nil {
		if now.Unix() >= job.next.Unix() {
			job.next = job.expr.Next(now)
			if !job.proc.IsRunning() {
				_job.Map.Store("last", int(now.Unix()))
			}
			if err := job.proc.Run(); err != nil {
				common.Logger.LogError("cmdsvr.taskJob.ExecHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
			}
		}
		return
	}

	statusRaw, err := ioutil.ReadFile(getTaskStatusFile(_job))
	if err != nil {
		// pass
	} else {
		statusStr := strings.TrimSpace(string(statusRaw))
		if len(statusStr) == 10 {
			status, err := strconv.Atoi(statusStr)
			if err != nil {
				// pass
			} else {
				_job.Map.Store("last", status)
				if int(now.Unix())-status > common.Config.ReportInterval {
					if err := job.proc.Kill(); err != nil {
						common.Logger.LogError("cmdsvr.taskJob.ExecHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
					}
					return
				}
			}
		}
	}

	if err := job.proc.Run(); err != nil {
		common.Logger.LogError("cmdsvr.taskJob.ExecHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
	}
}

func (job *taskJob) ExitHandler(_job *jobctl.Job) {
	if job.expr != nil {
		for i := 0; i < common.Config.StopTimeout*10; i++ {
			if job.proc.IsRunning() {
				time.Sleep(time.Millisecond * 100)
			} else {
				break
			}
		}
	}
	for {
		if err := job.proc.Signal(syscall.SIGTERM); err != nil {
			common.Logger.LogError("cmdsvr.taskJob.ExitHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
	for i := 0; i < common.Config.StopTimeout*10; i++ {
		if job.proc.IsRunning() {
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
	for {
		if err := job.proc.Kill(); err != nil {
			common.Logger.LogError("cmdsvr.taskJob.ExitHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
	for {
		if job.proc.IsRunning() {
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
	if err := os.Remove(getTaskStatusFile(_job)); err != nil {
		if !os.IsNotExist(err) {
			common.Logger.LogError("cmdsvr.taskJob.ExitHandler", "%v (%s)", err, common.JobManager.GetJobName(_job))
		}
	}
	common.JobManager.DestroyJob(_job)
}

func newTaskAddCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&taskAddCommand{Name: name, Argc: -1, AuthHandler: handler})
}
