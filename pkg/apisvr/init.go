package apisvr

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/daemon"
	"github.com/kimkit/appdaemon/static"
)

func init() {
	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	common.ApiSvr.Engine.Use(cors.New(config))

	common.ApiSvr.Register(&LoginController{Path: "/Login"})
	common.ApiSvr.Register(&GetLoginUserController{Path: "/GetLoginUser"})
	common.ApiSvr.Register(&GetLuaScriptListController{Path: "/GetLuaScriptList"})
	common.ApiSvr.Register(&GetLuaScriptController{Path: "/GetLuaScript"})
	common.ApiSvr.Register(&AddLuaScriptController{Path: "/AddLuaScript"})
	common.ApiSvr.Register(&UpdateLuaScriptController{Path: "/UpdateLuaScript"})
	common.ApiSvr.Register(&UpdateLuaScriptStatusController{Path: "/UpdateLuaScriptStatus"})
	common.ApiSvr.Register(&DeleteLuaScriptController{Path: "/DeleteLuaScript"})
	common.ApiSvr.Register(&GetServerListController{Path: "/GetServerList"})
	common.ApiSvr.Register(&UpdateServerStatusController{Path: "/UpdateServerStatus"})
	common.ApiSvr.Register(&DeleteServerController{Path: "/DeleteServer"})
	common.ApiSvr.Register(&GetTaskListController{Path: "/GetTaskList"})
	common.ApiSvr.Register(&GetTaskController{Path: "/GetTask"})
	common.ApiSvr.Register(&AddTaskController{Path: "/AddTask"})
	common.ApiSvr.Register(&UpdateTaskController{Path: "/UpdateTask"})
	common.ApiSvr.Register(&UpdateTaskStatusController{Path: "/UpdateTaskStatus"})
	common.ApiSvr.Register(&DeleteTaskController{Path: "/DeleteTask"})
	common.ApiSvr.Register(&GetOutputController{Path: "/GetOutput"})

	if handler, err := static.NewHandler("/"); err != nil {
		common.Logger.LogError("apisvr.init", "%v", err)
	} else {
		hf := gin.WrapH(handler)
		common.ApiSvr.Engine.GET("/favicon.ico", hf)
		common.ApiSvr.Engine.GET("/index.html", hf)
		common.ApiSvr.Engine.GET("/css/:file", hf)
		common.ApiSvr.Engine.GET("/js/:file", hf)
		common.ApiSvr.Engine.GET("/fonts/:file", hf)
		common.ApiSvr.Engine.GET("/img/:file", hf)
		common.ApiSvr.Engine.GET("/", hf)
	}
}

func Run() {
	if common.Config.Daemon {
		daemon.Daemon(common.Config.UI.LogFile, common.Config.UI.PidFile)
	}
	common.ApiSvr.ListenAndServe(common.Config.UI.Addr)
}
