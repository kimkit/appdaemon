package apisvr

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/luactl"
)

type LoginController struct {
	BaseController
	Path string
}

func (c *LoginController) POST(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" {
		c.Failure(ctx, ErrUsernameEmpty)
		return
	}
	if password == "" {
		c.Failure(ctx, ErrPasswordEmpty)
		return
	}

	got := false
	for u, p := range common.Config.UI.User {
		if u == username {
			if p != password {
				c.Failure(ctx, ErrPasswordInvalid)
				return
			}
			got = true
			break
		}
	}
	if !got {
		c.Failure(ctx, ErrUsernameNotExist)
		return
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	md5sum := luactl.Md5sum(fmt.Sprintf(
		"%s|%s|%s|%s",
		common.Config.UI.Salt,
		username,
		password,
		timestamp,
	))
	token := fmt.Sprintf("%s%s", md5sum, timestamp)
	c.Success(ctx, gin.H{
		"token": token,
	})
}
