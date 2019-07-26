package common

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kimkit/config"
	"github.com/kimkit/dbutil"
	"github.com/kimkit/ginsvr"
	"github.com/kimkit/jobext"
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
		Daemon         bool                    `json:"daemon"`
		LogFile        string                  `json:"-"`
		PidFile        string                  `json:"-"`
		Addr           string                  `json:"addr"`
		Port           int                     `json:"-"`
		Passwords      []string                `json:"passwords"`
		JobsFile       string                  `json:"-"`
		LogsDir        string                  `json:"logsdir"`
		ReportInterval int                     `json:"reportinterval"`
		StopTimeout    int                     `json:"stoptimeout"`
		Tasks          []*TaskConfig           `json:"tasks"`
		Http           map[string]*HttpConfig  `json:"http"`
		Redis          map[string]*RedisConfig `json:"redis"`
		DB             map[string]*DBConfig    `json:"db"`
		LuaScript      struct {
			Enable       bool              `json:"-"`
			DBConfig     *DBConfig         `json:"dbconfig"`
			Sql          string            `json:"sql"`
			IdName       string            `json:"idname"`
			IdInit       int               `json:"idinit"`
			FilterPrefix []string          `json:"filterprefix"`
			Libs         map[string]string `json:"libs"`
			Bootstrap    map[string]string `json:"bootstrap"`
		} `json:"luascript"`
		UI struct {
			Run     bool              `json:"-"`
			LogFile string            `json:"-"`
			PidFile string            `json:"-"`
			Addr    string            `json:"addr"`
			Salt    string            `json:"salt"`
			User    map[string]string `json:"user"`
		} `json:"ui"`
	}{}
	JobManager     = jobext.NewJobManager()
	ApiSvr         = ginsvr.NewServer()
	CmdSvr         = redsvr.NewServer()
	RedisClient    *redis.Client
	Logger         = logger.NewLogger()
	LuaScriptStore = luactl.NewLuaScriptStore(luactl.LuaScriptStoreOptions{CreateStateHandler: CreateStateHandler})
	HttpClient     = reqctl.NewClient(3)
	DBClient       *dbutil.DBWrapper
	Http           = make(map[string]*http.Client)
	RedisStore     sync.Map
	Redis          = make(map[string]*redis.Client)
	DB             = make(map[string]*dbutil.DBWrapper)
)

func init() {
	disableDaemon := flag.Bool("disable-daemon", false, "Disable daemon")
	enablePredefinedTasks := flag.Bool("enable-predefined-tasks", false, "Enable predefined tasks")
	enableLuaScript := flag.Bool("enable-lua-script", false, "Enable lua script")
	ui := flag.Bool("ui", false, "Run UI server")

	config.Load(&Config)
	if *disableDaemon {
		Config.Daemon = false
	}
	if !*enablePredefinedTasks {
		Config.Tasks = nil
	}
	Config.LuaScript.Enable = *enableLuaScript
	if !*enableLuaScript {
		Config.LuaScript.Bootstrap = nil
	}
	Config.UI.Run = *ui

	if Config.Addr == "" {
		Config.Addr = ":6380"
	}
	if !strings.Contains(Config.Addr, ":") {
		Config.Addr += ":6380"
	}
	arr := strings.Split(Config.Addr, ":")
	port, _ := strconv.Atoi(arr[len(arr)-1])
	Config.Port = port
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
	Config.UI.LogFile = fmt.Sprintf("%s%c%s", Config.LogsDir, os.PathSeparator, "appdaemon-ui.log")
	Config.UI.PidFile = fmt.Sprintf("%s%c%s", Config.LogsDir, os.PathSeparator, "appdaemon-ui.pid")
	if Config.UI.Addr == "" {
		Config.UI.Addr = ":7982"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("127.0.0.1:%s", strings.Split(Config.Addr, ":")[1]),
		Password: password,
	})
	initHttp()
	initRedis()
	initDB()
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

func initHttp() {
	for name, config := range Config.Http {
		if config.Alias == "" {
			Http[name] = NewHttp(config)
		}
	}
	for name, config := range Config.Http {
		if config.Alias != "" {
			if client, ok := Http[config.Alias]; ok {
				Http[name] = client
			}
		}
	}
	Http["#"] = HttpClient
}

func initRedis() {
	for name, config := range Config.Redis {
		if config.Alias == "" {
			Redis[name] = NewRedis(config)
		}
	}
	for name, config := range Config.Redis {
		if config.Alias != "" {
			if client, ok := Redis[config.Alias]; ok {
				Redis[name] = client
			}
		}
	}
	Redis["#"] = RedisClient
}

func initDB() {
	for name, config := range Config.DB {
		if config.Alias == "" {
			DB[name] = NewDB(config)
		}
	}
	for name, config := range Config.DB {
		if config.Alias != "" {
			if dbw, ok := DB[config.Alias]; ok {
				DB[name] = dbw
			}
		}
	}
	if Config.LuaScript.DBConfig != nil {
		alias := Config.LuaScript.DBConfig.Alias
		if alias == "" {
			DBClient = NewDB(Config.LuaScript.DBConfig)
			DB["#"] = DBClient
		} else {
			if dbw, ok := DB[alias]; ok {
				DBClient = dbw
				DB["#"] = dbw
			}
		}
	}
}

type HttpConfig struct {
	Alias   string `json:"alias"`
	Timeout int    `json:"timeout"`
}

func NewHttp(config *HttpConfig) *http.Client {
	return reqctl.NewClient(config.Timeout)
}

type RedisConfig struct {
	Alias        string `json:"alias"`
	Addr         string `json:"addr"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"poolsize"`
	MinIdleConns int    `json:"minidleconns"`
}

func NewRedis(config *RedisConfig) *redis.Client {
	key := fmt.Sprintf(
		"%s|%s|%d|%d|%d",
		config.Addr,
		config.Password,
		config.DB,
		config.PoolSize,
		config.MinIdleConns,
	)
	client, ok := RedisStore.Load(key)
	if !ok {
		client, _ = RedisStore.LoadOrStore(key, redis.NewClient(&redis.Options{
			Addr:         config.Addr,
			Password:     config.Password,
			DB:           config.DB,
			PoolSize:     config.PoolSize,
			MinIdleConns: config.MinIdleConns,
		}))
	}
	return client.(*redis.Client)
}

type DBConfig struct {
	Alias           string `json:"alias"`
	Driver          string `json:"driver"`
	DSN             string `json:"dsn"`
	MaxOpenConns    int    `json:"maxopenconns"`
	MaxIdleConns    int    `json:"maxidleconns"`
	ConnMaxLifetime int    `json:"connmaxlifetime"`
}

func NewDB(config *DBConfig) *dbutil.DBWrapper {
	if config.Driver == "" {
		config.Driver = "mysql"
	}
	return dbutil.NewDB(
		config.Driver,
		config.DSN,
		config.MaxOpenConns,
		config.MaxIdleConns,
		config.ConnMaxLifetime,
	)
}

var (
	ErrHttpNotFound  = fmt.Errorf("ErrHttpNotFound")
	ErrRedisNotFound = fmt.Errorf("ErrRedisNotFound")
	ErrDBNotFound    = fmt.Errorf("ErrDBNotFound")
)

func GetHttp(name string) (*http.Client, error) {
	client, ok := Http[name]
	if !ok {
		return nil, ErrHttpNotFound
	}
	return client, nil
}

func GetRedis(name string) (*redis.Client, error) {
	client, ok := Redis[name]
	if !ok {
		return nil, ErrRedisNotFound
	}
	return client, nil
}

func GetDB(name string) (*dbutil.DBWrapper, error) {
	dbw, ok := DB[name]
	if !ok {
		return nil, ErrDBNotFound
	}
	return dbw, nil
}
