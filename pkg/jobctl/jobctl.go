package jobctl

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrJobStoppingTimeout = fmt.Errorf("job stopping timetout")
)

type JobHandler func(*Job)

type Job struct {
	status      int64
	Map         sync.Map
	InitHandler JobHandler
	ExecHandler JobHandler
	ExitHandler JobHandler
}

func (job *Job) IsRunning() bool {
	return !atomic.CompareAndSwapInt64(&job.status, 0, 0)
}

func (job *Job) Stop(waitTime int) error {
	if !atomic.CompareAndSwapInt64(&job.status, 1, 2) {
		if job.IsRunning() {
			goto check
		} else {
			return nil
		}
	}

check:

	if waitTime <= 0 {
		return nil
	}

	for i := 0; i < waitTime*10; i++ {
		if !job.IsRunning() {
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}

	return ErrJobStoppingTimeout
}

func (job *Job) Start() error {
	//  <---------
	// /          \
	// 0 --> 1 --> 2

	if !atomic.CompareAndSwapInt64(&job.status, 0, 1) {
		return nil
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				debug.PrintStack()
				job.Stop(0)
			}
			atomic.CompareAndSwapInt64(&job.status, 2, 0)
		}()
		if job.ExitHandler != nil {
			defer job.ExitHandler(job)
		}
		if job.InitHandler != nil {
			job.InitHandler(job)
		}
		for {
			if atomic.CompareAndSwapInt64(&job.status, 2, 2) {
				return
			}
			if job.ExecHandler == nil {
				job.Stop(0)
				return
			}
			job.ExecHandler(job)
		}
	}()

	return nil
}

func NewJob(ptr interface{}) *Job {
	job := &Job{}
	jobVal := reflect.ValueOf(job)
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Kind() != reflect.Ptr {
		return nil
	}
	if ptrVal.Elem().Kind() != reflect.Struct {
		return nil
	}
	for _, handlerName := range []string{"InitHandler", "ExecHandler", "ExitHandler"} {
		handlerVal := ptrVal.MethodByName(handlerName)
		if handlerVal.IsValid() && handlerVal.Type().AssignableTo(reflect.TypeOf(job.InitHandler)) {
			jobVal.Elem().FieldByName(handlerName).Set(handlerVal)
		}
	}
	return job
}
