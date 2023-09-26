package cmdlib

import (
	"github.com/kimkit/appdaemon/pkg/redsvr"
)

type pingCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *pingCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	reply := "PONG"
	if len(args) == 1 {
		reply = args[0]
	}
	redsvr.WriteSimpleString(conn, reply)
	return nil
}

func NewPingCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&pingCommand{Name: name, Argc: -1, AuthHandler: handler})
}
