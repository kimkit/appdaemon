package cmdsvr

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/redsvr"
	"github.com/kimkit/appdaemon/pkg/thread"
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

	fp, err := os.Open(common.Config.LogsDir)
	if err != nil {
		return err
	}
	fis, err := fp.Readdir(0)
	if err != nil {
		return err
	}
	var files []string
	for _, fi := range fis {
		file := fi.Name()
		if strings.HasSuffix(file, ".output") && (file == args[0]+".output" || strings.HasPrefix(file, args[0]+"_")) {
			files = append(files, file)
		}
	}
	if len(files) == 0 {
		return fmt.Errorf("no matched files")
	}

	tfs := make(map[string]*tail.Tail)
	for _, file := range files {
		tf, err := tail.TailFile(path.Join(common.Config.LogsDir, file), tail.Config{
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: os.SEEK_END,
			},
			Follow: true,
			Poll:   true,
			ReOpen: true,
			Logger: tail.DiscardingLogger,
		})
		if err != nil {
			return err
		}
		tfs[file] = tf
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

	tm := thread.ThreadManager{}
	exit := make(chan error)

	for file, tf := range tfs {
		tm.New(func(file string, tf *tail.Tail) func() {
			return func() {
				for {
					select {
					case line, ok := <-tf.Lines:
						if !ok {
							common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "file tail stopped (%s)", file)
							return
						}
						if line.Err != nil {
							common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "%v (%s)", line.Err, file)
						} else {
							err := redsvr.WriteArray(conn, []interface{}{
								"message",
								args[0],
								fmt.Sprintf("%s: %s", file, line.Text),
							})
							if err != nil {
								common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "%v (%s)", err, file)
								return
							}
						}
					case <-exit:
						common.Logger.LogInfo("cmdsvr.subscribeCommand.S1Handler", "tail exit (%s)", file)
						tf.Stop()
						tf.Cleanup()
						return
					}
				}
			}
		}(file, tf))
	}

	tm.New(func() {
		buf := make([]byte, 1)
		for {
			if _, err := conn.Read(buf); err != nil {
				common.Logger.LogError("cmdsvr.subscribeCommand.S1Handler", "reader stopped: %v (%s)", err, args[0])
				close(exit)
				return
			}
		}
	})

	tm.Wait()
	common.Logger.LogInfo("cmdsvr.subscribeCommand.S1Handler", "sender stopped (%s)", args[0])
	return nil
}

func newSubscribeCommand(name string, handler redsvr.CommandHandler) *redsvr.Command {
	return redsvr.NewCommand(&subscribeCommand{Name: name, Argc: 1, AuthHandler: handler})
}
