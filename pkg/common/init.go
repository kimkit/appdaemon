package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/kimkit/config"
	"github.com/kimkit/jobext"
	"github.com/kimkit/redsvr"
)

var (
	Config = struct {
		Daemon         bool     `json:"daemon"`
		Addr           string   `json:"addr"`
		Passwords      []string `json:"passwords"`
		JobsFile       string   `json:"jobsfile"`
		LogsDir        string   `json:"logsdir"`
		ReportInterval int      `json:"reportinterval"`
		StopTimeout    int      `json:"stoptimeout"`
	}{}
	JobManager = jobext.NewJobManager()
	Cmdsvr     = redsvr.NewServer()
	Client     *redis.Client
)

func init() {
	config.Load(&Config)
	if Config.Addr == "" {
		Config.Addr = ":6380"
	}
	if !strings.Contains(Config.Addr, ":") {
		Config.Addr += ":6380"
	}
	if Config.JobsFile == "" {
		Config.JobsFile = ".jobs.json"
	}
	if Config.LogsDir == "" {
		Config.LogsDir = ".logs"
	}
	if err := os.MkdirAll(Config.LogsDir, 0755); err != nil {
		redsvr.Log("ERROR", "common.init: %v (%s)", err, Config.LogsDir)
	}
	if Config.ReportInterval <= 0 {
		Config.ReportInterval = 10
	}
	if Config.StopTimeout <= 0 {
		Config.StopTimeout = 10
	}
	password := ""
	if len(Config.Passwords) > 0 {
		password = Config.Passwords[0]
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("127.0.0.1:%s", strings.Split(Config.Addr, ":")[1]),
		Password: password,
	})
}
