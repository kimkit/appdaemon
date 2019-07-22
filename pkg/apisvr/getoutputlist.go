package apisvr

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/dbutil"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type GetOutputListController struct {
	BaseController
	Path string
}

func (c *GetOutputListController) GET(ctx *gin.Context) {
	if err := c.CheckPermission(ctx); err != nil {
		// c.Failure(ctx, err)
		// return
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
	}

	defer ws.Close()

	id := "0"
	limit := 100

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
				return
			}
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if len(rows) == 0 {
			if err := ws.WriteJSON(gin.H{"code": -2, "message": "empty"}); err != nil {
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
				return
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}
