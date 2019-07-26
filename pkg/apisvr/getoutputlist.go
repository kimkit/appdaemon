package apisvr

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
)

type GetOutputListController struct {
	BaseController
	Path string
}

func (c *GetOutputListController) GET(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
		c.Failure(ctx, err)
		return
	}

	name := ctx.Query("name")
	if name == "" {
		c.Failure(ctx, ErrLuaScriptNameEmpty)
		return
	}

	db, err := common.GetDB("#")
	if err != nil {
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
			common.Logger.LogDebug("apisvr.GetOutputListController.GET", "%v", err)
		}
		common.Logger.LogDebug("apisvr.GetOutputListController.GET", "exit")
	}()

	go func() {
		defer func() {
			common.Logger.LogDebug("apisvr.GetOutputListController.GET", "exit reading")
		}()
		for {
			typ, msg, err := ws.ReadMessage()
			if err != nil {
				common.Logger.LogDebug("apisvr.GetOutputListController.GET", "%v", err)
				return
			}
			common.Logger.LogDebug("apisvr.GetOutputListController.GET", "%v %v", typ, string(msg))
			if typ == websocket.PingMessage {
				if err := ws.WriteMessage(websocket.PongMessage, []byte("pong")); err != nil {
					common.Logger.LogDebug("apisvr.GetOutputListController.GET", "%v", err)
				}
			}
		}
	}()

	id := "0"
	limitStr := ctx.Query("limit")
	limit := 0
	if limitStr == "" {
		limit = 100
	} else {
		_limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 100
		} else {
			if _limit <= 0 {
				limit = 100
			} else if _limit > 500 {
				limit = 500
			} else {
				limit = _limit
			}
		}
	}

	for {
		var sql string
		var first bool
		if id == "0" {
			first = true
			sql = fmt.Sprintf("select id,addr,line from output where name = ? order by id desc limit %d", limit)
		} else {
			first = false
			sql = fmt.Sprintf("select id,addr,line from output where name = ? and id > %s order by id asc limit %d", id, limit)
		}
		rows, err := dbutil.FetchAll(db.Query(sql, name))
		if err != nil {
			if err := ws.WriteJSON(gin.H{"code": -1, "message": err.Error()}); err != nil {
				common.Logger.LogDebug("apisvr.GetOutputListController.GET", "%v", err)
				return
			}
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if len(rows) == 0 {
			if err := ws.WriteJSON(gin.H{"code": -2, "message": "empty"}); err != nil {
				common.Logger.LogDebug("apisvr.GetOutputListController.GET", "%v", err)
				return
			}
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if first {
			var _rows []map[string]string
			for i := len(rows) - 1; i >= 0; i-- {
				_rows = append(_rows, rows[i])
			}
			rows = _rows
		}
		for _, row := range rows {
			id = row["id"]
			if err := ws.WriteJSON(gin.H{"code": 1, "message": "OK", "data": row}); err != nil {
				common.Logger.LogDebug("apisvr.GetOutputListController.GET", "%v", err)
				return
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}
