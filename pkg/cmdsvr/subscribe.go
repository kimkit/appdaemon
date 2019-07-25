package cmdsvr

import (
	"os"
	"path"

	"github.com/hpcloud/tail"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/redsvr"
)

type subscribeCommand struct {
	Name        string
	Argc        int
	AuthHandler redsvr.CommandHandler
}

func (cmd *subscribeCommand) S1Handler(_cmd *redsvr.Command, args []string, conn *redsvr.Conn) error {
	if cmd.AuthHandler != nil {
		if err := cmd.AuthHandler(_cmd, args, conn); err != nil {
			return err
		}
	}
	file := path.Join(common.Config.LogsDir, args[0]+".output")
	if _, err := os.Stat(file); err != nil {
		return err
	}

	tf, err := tail.TailFile(file, tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
		Follow: true,
		Logger: tail.DiscardingLogger,
	})
	if err != nil {
		return err
	}

	err = redsvr.WriteArray(conn, []interface{}{
		"subscribe",
		args[0],
		1,
	})
	if err != nil {
		common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "%v (%s)", err, args[0])
		return nil
	}

	errch := make(chan error, 1)
	go func() {
		defer common.Logger.LogInfo("cmdsvr.subscribeCommand.S1Handler", "reader stopped (%s)", args[0])
		buf := make([]byte, 1)
		for {
			if _, err := conn.Read(buf); err != nil {
				errch <- err
				return
			}
		}
	}()

	defer func() {
		conn.Close()
		common.Logger.LogInfo("cmdsvr.subscribeCommand.S1Handler", "sender stopped (%s)", args[0])
	}()

	for {
		select {
		case line, ok := <-tf.Lines:
			if !ok {
				common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "file tail stopped (%s)", args[0])
				return nil
			}
			if line.Err != nil {
				common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "%v (%s)", line.Err, args[0])
			} else {
				err := redsvr.WriteArray(conn, []interface{}{
					"message",
					args[0],
					line.Text,
				})
				if err != nil {
					common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "%v (%s)", err, args[0])
					return nil
				}
			}
		case err := <-errch:
			common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "%v (%s)", err, args[0])
			return nil
		}
	}
	return nil
}

func newSubscribeCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&subscribeCommand{Name: name, Argc: 1, AuthHandler: handler})
}
