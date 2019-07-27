package apisvr

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kimkit/apires"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/luactl"
)

var (
	ErrSystemError                = apires.NewError(-1000, "ErrSystemError", nil)
	ErrUserNotLogin               = apires.NewError(-1001, "ErrUserNotLogin", nil)
	ErrTokenEmpty                 = apires.NewError(-1002, "ErrTokenEmpty", nil)
	ErrTokenFormatInvalid         = apires.NewError(-1003, "ErrTokenFormatInvalid", nil)
	ErrTokenInvalid               = apires.NewError(-1004, "ErrTokenInvalid", nil)
	ErrTokenExpired               = apires.NewError(-1005, "ErrTokenExpired", nil)
	ErrUsernameEmpty              = apires.NewError(-1006, "ErrUsernameEmpty", nil)
	ErrPasswordEmpty              = apires.NewError(-1007, "ErrPasswordEmpty", nil)
	ErrPasswordInvalid            = apires.NewError(-1008, "ErrPasswordInvalid", nil)
	ErrUsernameNotExist           = apires.NewError(-1009, "ErrUsernameNotExist", nil)
	ErrLuaScriptNotExist          = apires.NewError(-1010, "ErrLuaScriptNotExist", nil)
	ErrLuaScriptNameEmpty         = apires.NewError(-1011, "ErrLuaScriptNameEmpty", nil)
	ErrLuaScriptNameFormatInvalid = apires.NewError(-1012, "ErrLuaScriptNameFormatInvalid", nil)
	ErrLuaScriptEmpty             = apires.NewError(-1013, "ErrLuaScriptEmpty", nil)
	ErrLuaScriptSyntaxError       = apires.NewError(-1014, "ErrLuaScriptSyntaxError", nil)
	ErrLuaScriptExist             = apires.NewError(-1015, "ErrLuaScriptExist", nil)
	ErrLuaScriptNameDuplicate     = apires.NewError(-1016, "ErrLuaScriptNameDuplicate", nil)
	ErrServerAddrNotExist         = apires.NewError(-1017, "ErrServerAddrNotExist", nil)
	ErrTaskNameEmpty              = apires.NewError(-1018, "ErrTaskNameEmpty", nil)
	ErrTaskNameFormatInvalid      = apires.NewError(-1019, "ErrTaskNameFormatInvalid", nil)
	ErrTaskRuleEmpty              = apires.NewError(-1020, "ErrTaskRuleEmpty", nil)
	ErrTaskRuleInvalid            = apires.NewError(-1021, "ErrTaskRuleInvalid", nil)
	ErrTaskCommandEmpty           = apires.NewError(-1022, "ErrTaskCommandEmpty", nil)
	ErrTaskCommandInvalid         = apires.NewError(-1023, "ErrTaskCommandInvalid", nil)
	ErrTaskExist                  = apires.NewError(-1024, "ErrTaskExist", nil)
	ErrTaskNameDuplicate          = apires.NewError(-1025, "ErrTaskNameDuplicate", nil)
	ErrJobNameEmpty               = apires.NewError(-1026, "ErrJobNameEmpty", nil)
	ErrServerAddrEmpty            = apires.NewError(-1027, "ErrServerAddrEmpty", nil)
	ErrJobIsRunning               = apires.NewError(-1028, "ErrJobIsRunning", nil)
	ErrTaskNotExist               = apires.NewError(-1029, "ErrTaskNotExist", nil)
	ErrTaskEnable                 = apires.NewError(-1030, "ErrTaskEnable", nil)
	ErrLuaScriptEnable            = apires.NewError(-1031, "ErrLuaScriptEnable", nil)
)

type BaseController struct{}

func (c *BaseController) Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"code":    1,
		"message": "OK",
		"data":    data,
	})
}

func (c *BaseController) Failure(ctx *gin.Context, err error) {
	if _err, ok := err.(*apires.Error); ok {
		ctx.JSON(200, gin.H{
			"code":    _err.Reply.Code,
			"message": _err.Reply.Message,
			"data":    _err.Reply.Data,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    ErrSystemError.(*apires.Error).Reply.Code,
			"message": ErrSystemError.(*apires.Error).Reply.Message + ": " + err.Error(),
			"data":    ErrSystemError.(*apires.Error).Reply.Data,
		})
	}
}

func (c *BaseController) GetLoginUser(ctx *gin.Context) (string, error) {
	token := ctx.Query("token")
	if token == "" {
		return "", ErrTokenEmpty
	}
	if len(token) != 42 { // md5sum + timestamp
		return "", ErrTokenFormatInvalid
	}

	md5sum := token[0:32]
	timestamp := token[32:]

	user := ""
	for username, password := range common.Config.UI.User {
		if luactl.Md5sum(fmt.Sprintf("%s|%s|%s|%s", common.Config.UI.Salt, username, password, timestamp)) == md5sum {
			user = username
			break
		}
	}
	if user == "" {
		return "", ErrTokenInvalid
	}

	t, err := strconv.Atoi(timestamp)
	if err != nil {
		return "", err
	}
	if time.Now().Unix()-int64(t) > 7*3600*24 { // 7 days
		return "", ErrTokenExpired
	}

	return user, nil
}

func (c *BaseController) CheckPermission(ctx *gin.Context) error {
	if _, err := c.GetLoginUser(ctx); err != nil {
		return err
	}
	return nil
}
