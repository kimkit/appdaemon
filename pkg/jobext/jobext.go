package jobext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/kimkit/appdaemon/pkg/jobctl"
)

type JobManager struct {
	Map sync.Map
}

func NewJobManager() *JobManager {
	return &JobManager{}
}

func (jm *JobManager) GetJob(name string, ptr interface{}) *jobctl.Job {
	if ptr == nil {
		jobVal, _ := jm.Map.Load(name)
		job, _ := jobVal.(*jobctl.Job)
		return job
	} else {
		var job *jobctl.Job
		if _job, ok := ptr.(*jobctl.Job); ok {
			job = _job
		} else {
			job = jobctl.NewJob(ptr)
		}
		if job == nil {
			return nil
		} else {
			job.Map.Store("_job_name", name)
			jobVal, _ := jm.Map.LoadOrStore(name, job)
			job, _ := jobVal.(*jobctl.Job)
			return job
		}
	}
}

func (jm *JobManager) GetJobName(job *jobctl.Job) string {
	if job == nil {
		return ""
	}
	nameVal, _ := job.Map.Load("_job_name")
	name, _ := nameVal.(string)
	return name
}

func (jm *JobManager) DestroyJob(job *jobctl.Job) {
	if job != nil {
		if nameVal, ok := job.Map.Load("_job_name"); ok {
			jm.Map.Delete(nameVal)
		}
	}
}

func (jm *JobManager) SetJobInfo(job *jobctl.Job, v ...interface{}) {
	if job != nil {
		job.Map.Store("_job_info", v)
	}
}

func (jm *JobManager) GetJobInfo(job *jobctl.Job) []interface{} {
	if job != nil {
		infoVal, _ := job.Map.Load("_job_info")
		info, _ := infoVal.([]interface{})
		return info
	}
	return nil
}

func (jm *JobManager) StopJob(job *jobctl.Job) {
	for {
		if err := job.Stop(10); err != nil {
			log.Printf("ERROR jobext.JobManager.StopJob: %v (%s)", err, jm.GetJobName(job))
		} else {
			break
		}
	}
}

func (jm *JobManager) StopAllJobs() {
	var wg sync.WaitGroup
	jm.Map.Range(func(k, v interface{}) bool {
		job, _ := v.(*jobctl.Job)
		if job != nil {
			wg.Add(1)
			go func(job *jobctl.Job) {
				defer wg.Done()
				jm.StopJob(job)
			}(job)
		}
		return true
	})
	wg.Wait()
}

func (jm *JobManager) GetRunningJobs() []string {
	var list []string
	jm.Map.Range(func(k, v interface{}) bool {
		job, _ := v.(*jobctl.Job)
		if job != nil {
			if job.IsRunning() {
				if name, ok := k.(string); ok {
					list = append(list, name)
				}
			}
		}
		return true
	})
	sort.Strings(list)
	return list
}

func (jm *JobManager) SaveRunningJobs(file string) {
	list := "[\n"
	jm.Map.Range(func(k, v interface{}) bool {
		job, _ := v.(*jobctl.Job)
		if job != nil {
			if job.IsRunning() {
				if info := jm.GetJobInfo(job); len(info) > 0 {
					if infoBytes, err := json.Marshal(info); err != nil {
						log.Printf("ERROR jobext.JobManager.SaveRunningJobs: %v (%s)", err, jm.GetJobName(job))
					} else {
						list += fmt.Sprintf("  %s,\n", string(infoBytes))
					}
				}
			}
		}
		return true
	})
	if len(list) > 2 {
		list = list[0:len(list)-2] + "\n]\n"
	} else {
		list = "[]"
	}
	if err := ioutil.WriteFile(file, []byte(list), 0644); err != nil {
		log.Printf("ERROR jobext.JobManager.SaveRunningJobs: %v (%s)", err, file)
	}
}

type JobLoader func([]interface{}) error

func (jm *JobManager) LoadJobs(file string, loader JobLoader, infos ...[]interface{}) {
	if loader == nil {
		return
	}
	listBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("ERROR jobext.JobManager.LoadJobs: %v (%s)", err, file)
		if len(infos) == 0 {
			return
		}
	}
	var list [][]interface{}
	if err := json.Unmarshal(listBytes, &list); err != nil {
		log.Printf("ERROR jobext.JobManager.LoadJobs: %v (%s)", err, file)
		if len(infos) == 0 {
			return
		}
	}

	for _, info := range infos {
		list = append(list, info)
	}

	retry := 0
	job := &jobctl.Job{}
	job.InitHandler = func(job *jobctl.Job) {
		time.Sleep(time.Second)
	}
	job.ExecHandler = func(job *jobctl.Job) {
		if len(list) == 0 {
			job.Stop(0)
			return
		}
		if err := loader(list[0]); err != nil {
			log.Printf("ERROR jobext.JobManager.LoadJobs: %v (%v)", err, list[0])
			retry++
			if retry >= 5 {
				list = list[1:]
				retry = 0
			}
			time.Sleep(time.Second)
			return
		}
		list = list[1:]
		retry = 0
	}
	job.ExitHandler = func(job *jobctl.Job) {
		jm.DestroyJob(job)
	}
	jm.GetJob("_job_loader", job)
	if err := job.Start(); err != nil {
		log.Printf("ERROR jobext.JobManager.LoadJobs: %v (%s)", err, jm.GetJobName(job))
	}
}

func (jm *JobManager) RunJob(job *jobctl.Job, action string) (int, error) {
	if job == nil {
		return 0, fmt.Errorf("job nil")
	}
	switch action {
	case "start":
		if err := job.Start(); err != nil {
			return 0, err
		} else {
			return 1, nil
		}
	case "stop":
		if err := job.Stop(10); err != nil {
			return 0, err
		} else {
			return 1, nil
		}
	case "status":
		if !job.IsRunning() {
			return 0, nil
		} else {
			return 1, nil
		}
	default:
		return 0, fmt.Errorf("action invalid")
	}
}
