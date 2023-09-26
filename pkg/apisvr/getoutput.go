package apisvr

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/kimkit/appdaemon/pkg/apires"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/dbutil"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type GetOutputController struct {
	BaseController
	Path string
}

func (c *GetOutputController) wsSuccess(conn *websocket.Conn, data interface{}) error {
	return conn.WriteJSON(gin.H{
		"code":    1,
		"message": "OK",
		"data":    data,
	})
}

func (c *GetOutputController) wsFailure(conn *websocket.Conn, err error) error {
	if _err, ok := err.(*apires.Error); ok {
		return conn.WriteJSON(gin.H{
			"code":    _err.Reply.Code,
			"message": _err.Reply.Message,
			"data":    _err.Reply.Data,
		})
	} else {
		return conn.WriteJSON(gin.H{
			"code":    ErrSystemError.(*apires.Error).Reply.Code,
			"message": ErrSystemError.(*apires.Error).Reply.Message + ": " + err.Error(),
			"data":    ErrSystemError.(*apires.Error).Reply.Data,
		})
	}
}

func (c *GetOutputController) GET(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
		c.Failure(ctx, err)
		return
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		c.Failure(ctx, err)
		return
	}

	defer func() {
		if err := ws.Close(); err != nil {
			common.Logger.LogError("apisvr.GetOutputController.GET", "%v", err)
		}
		common.Logger.LogInfo("apisvr.GetOutputController.GET", "exit")
	}()

	jobname := strings.TrimSpace(ctx.Query("jobname"))
	subaddr := strings.TrimSpace(ctx.Query("subaddr"))
	if jobname == "" {
		c.wsFailure(ws, ErrJobNameEmpty)
		return
	}
	if subaddr == "" {
		c.wsFailure(ws, ErrServerAddrEmpty)
		return
	}

	db, err := common.GetDB("#")
	if err != nil {
		c.wsFailure(ws, err)
		return
	}

	rows, err := dbutil.FetchAll(db.Query(
		"select id from server where addr = ? and status = 1",
		subaddr,
	))
	if err != nil {
		c.wsFailure(ws, err)
		return
	}
	if len(rows) == 0 {
		c.wsFailure(ws, ErrServerAddrNotExist)
		return
	}

	password := ""
	if len(common.Config.Passwords) > 0 {
		password = common.Config.Passwords[0]
	}
	client := common.NewRedis(&common.RedisConfig{
		Addr:     subaddr,
		Password: password,
	})
	pubsub := client.Subscribe(jobname)

	go func() {
		defer func() {
			if err := pubsub.Close(); err != nil {
				common.Logger.LogError("apisvr.GetOutputController.GET", "%v", err)
			}
			common.Logger.LogInfo("apisvr.GetOutputController.GET", "exit reading")
		}()
		for {
			typ, msg, err := ws.ReadMessage()
			if err != nil {
				common.Logger.LogError("apisvr.GetOutputController.GET", "%v", err)
				return
			}
			common.Logger.LogDebug("apisvr.GetOutputController.GET", "%v %s", typ, msg)
			if typ == websocket.PingMessage {
				if err := ws.WriteMessage(websocket.PongMessage, []byte("pong")); err != nil {
					common.Logger.LogError("apisvr.GetOutputController.GET", "%v", err)
					return
				}
			}
		}
	}()

	for {
		val, err := pubsub.ReceiveTimeout(time.Second)
		if err != nil {
			if !strings.Contains(err.Error(), "i/o timeout") {
				time.Sleep(time.Second)
				if err := c.wsSuccess(ws, err.Error()); err != nil {
					common.Logger.LogError("apisvr.GetOutputController.GET", "%v", err)
					return
				}
			}
		} else {
			if msg, ok := val.(*redis.Message); ok {
				if err := c.wsSuccess(ws, msg.Payload); err != nil {
					common.Logger.LogError("apisvr.GetOutputController.GET", "%v", err)
					return
				}
			}
		}
	}
}
