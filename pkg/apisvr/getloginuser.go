package apisvr

import (
	"github.com/gin-gonic/gin"
)

type GetLoginUserController struct {
	BaseController
	Path string
}

func (c *GetLoginUserController) GET(ctx *gin.Context) {
	user, err := c.GetLoginUser(ctx)
	if err != nil {
		c.Failure(ctx, err)
		return
	}
	c.Success(ctx, gin.H{
		"user": user,
	})
}
