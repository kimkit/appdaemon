package cmdlib

import (
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type echoCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *echoCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	redsvr.WriteBlukString(conn, args[0])
	return nil
}

func NewEchoCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&echoCommand{Name: name, Argc: 1, AuthHandler: handler})
}
