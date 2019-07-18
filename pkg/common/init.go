package common

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kimkit/config"
	"github.com/kimkit/jobext"
	"github.com/kimkit/lister"
	"github.com/kimkit/logger"
	"github.com/kimkit/luactl"
	"github.com/kimkit/redsvr"
	"github.com/kimkit/reqctl"
)

type TaskConfig struct {
	Name string   `json:"name"`
	Rule string   `json:"rule"`
	Args []string `json:"args"`
}

var (
	Config = struct {
		Daemon         bool          `json:"daemon"`
		LogFile        string        `json:"-"`
		PidFile        string        `json:"-"`
		Addr           string        `json:"addr"`
		Passwords      []string      `json:"passwords"`
		JobsFile       string        `json:"-"`
		LogsDir        string        `json:"logsdir"`
		ReportInterval int           `json:"reportinterval"`
		StopTimeout    int           `json:"stoptimeout"`
		Tasks          []*TaskConfig `json:"tasks"`
		Dsn            string        `json:"dsn"`
		Sql            string        `json:"sql"`
		IdName         string        `json:"idname"`
		IdInit         int           `json:"idinit"`
	}{}
	JobManager     = jobext.NewJobManager()
	Cmdsvr         = redsvr.NewServer()
	RedisClient    *redis.Client
	Logger         = logger.NewLogger()
	Lister         lister.Lister
	LuaScriptStore = luactl.NewLuaScriptStore(luactl.LuaScriptStoreOptions{CreateStateHandler: CreateStateHandler})
	HttpClient     = reqctl.NewClient(10)
)

func init() {
	disableDaemon := flag.Bool("disable-daemon", false, "Disable daemon")
	enablePredefinedTasks := flag.Bool("enable-predefined-tasks", false, "Enable predefined tasks")
	enableLuaScript := flag.Bool("enable-lua-script", false, "Enable lua script")

	config.Load(&Config)
	if *disableDaemon {
		Config.Daemon = false
	}
	if !*enablePredefinedTasks {
		Config.Tasks = nil
	}
	if !*enableLuaScript {
		Config.Dsn = ""
	}

	if Config.Addr == "" {
		Config.Addr = ":6380"
	}
	if !strings.Contains(Config.Addr, ":") {
		Config.Addr += ":6380"
	}
	if Config.LogsDir == "" {
		Config.LogsDir = "logs"
	}
	if err := os.MkdirAll(Config.LogsDir, 0755); err != nil {
		Logger.LogError("common.init", "%v (%s)", err, Config.LogsDir)
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
	Config.LogFile = fmt.Sprintf("%s%c%s", Config.LogsDir, os.PathSeparator, "appdaemon.log")
	Config.PidFile = fmt.Sprintf("%s%c%s", Config.LogsDir, os.PathSeparator, "appdaemon.pid")
	Config.JobsFile = fmt.Sprintf("%s%c%s", Config.LogsDir, os.PathSeparator, "jobs.json")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("127.0.0.1:%s", strings.Split(Config.Addr, ":")[1]),
		Password: password,
	})
	if Config.Dsn != "" {
		db, err := sql.Open("mysql", Config.Dsn)
		if err != nil {
			Logger.LogError("common.init", "%v", err)
		} else {
			if Config.Sql != "" {
				if Config.IdName == "" {
					Config.IdName = "id"
				}
				Lister = lister.NewDBLister(db, Config.Sql, Config.IdName, Config.IdInit)
			}
		}
	}
	initLuaLib()
}

func GetTaskInfos() [][]interface{} {
	var infos [][]interface{}
	for _, task := range Config.Tasks {
		num, err := strconv.Atoi(task.Rule)
		if err != nil {
			var info []interface{}
			info = append(info, "task.add")
			info = append(info, task.Name)
			info = append(info, task.Rule)
			for _, arg := range task.Args {
				info = append(info, arg)
			}
			infos = append(infos, info)
		} else {
			for i := 0; i < num; i++ {
				var info []interface{}
				info = append(info, "task.add")
				info = append(info, fmt.Sprintf("%s_%03d", task.Name, i))
				info = append(info, "")
				for _, arg := range task.Args {
					info = append(info, arg)
				}
				infos = append(infos, info)
			}
		}
	}
	return infos
}
