package cmdsvr

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/jobctl"
	"github.com/kimkit/redsvr"
)

type taskDeleteCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *taskDeleteCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	key := GetTaskKey(args[0])
	var jobs []*jobctl.Job
	common.JobManager.Map.Range(func(k, v interface{}) bool {
		name, _ := k.(string)
		job, _ := v.(*jobctl.Job)
		if job != nil {
			if name == key {
				job.Stop(0)
				jobs = append(jobs, job)
				return true
			}
			if strings.HasPrefix(name, key+"_") {
				if _, err := strconv.Atoi(name[len(key)+1:]); err != nil {
					// ignore
				} else {
					job.Stop(0)
					jobs = append(jobs, job)
					return true
				}
			}
		}
		return true
	})

	if len(jobs) > 0 && len(args) > 1 {
		t, _ := strconv.Atoi(args[1])
		if t > 0 {
			done := false
			for i := 0; i < t*10; i++ {
				time.Sleep(time.Millisecond * 100)
				got := false
				for _, job := range jobs {
					if job.IsRunning() {
						got = true
						break
					}
				}
				if !got {
					done = true
					break
				}
			}
			if !done {
				return fmt.Errorf("wait job `%s_*` stopping timeout", key)
			}
		}
	}

	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

func newTaskDeleteCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&taskDeleteCommand{Name: name, Argc: -1, AuthHandler: handler})
}
