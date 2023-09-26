package cmdlib

import (
	"fmt"
	"strings"

	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type authCommand struct {
	Name string
	Argc int
	Pwds []string
}

func (cmd *authCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if len(cmd.Pwds) > 0 {
		pass := false
		for _, pwd := range cmd.Pwds {
			if pwd == args[0] {
				pass = true
				conn.Map.Store("authenticated", true)
				conn.Map.Store("password", args[0])
				break
			}
		}
		if !pass {
			return fmt.Errorf("invalid password")
		}
	}
	redsvr.WriteSimpleString(conn, "OK")
	return nil
}

func NewAuthCommand(name string, pwds []string) *redsvr.Command {
	return redsvr.NewCommand(&authCommand{Name: name, Argc: 1, Pwds: pwds})
}

func CheckAuth(conn *redsvr.Conn) error {
	authenticatedVal, _ := conn.Map.Load("authenticated")
	authenticated, _ := authenticatedVal.(bool)
	if !authenticated {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func GetAuthUser(conn *redsvr.Conn) string {
	passwordVal, _ := conn.Map.Load("password")
	password, _ := passwordVal.(string)
	info := strings.Split(password, ":")
	if len(info) < 2 {
		return ""
	} else {
		return info[0]
	}
}
