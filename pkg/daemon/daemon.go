package daemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func Daemon(logfile, pidfile string) {
	if os.Getenv("__DAEMON__") != "true" {
		var lfp *os.File
		if logfile != "" {
			if fp, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
				log.Printf("WARNING daemon.Daemon(): open log file `%s` failed (%v)", logfile, err)
			} else {
				lfp = fp
			}
		}
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), "__DAEMON__=true")
		cmd.Stdin = nil
		cmd.Stdout = lfp
		cmd.Stderr = lfp
		if err := cmd.Start(); err != nil {
			log.Fatalf("ERROR daemon.Daemon(): create process failed (%v)", err)
		}
		os.Exit(0)
	}

	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(pidfile, []byte(pid), 0644); err != nil {
		log.Printf("WARNING daemon.Daemon(): write pid `%s` to file `%s` failed (%v)", pid, pidfile, err)
	}
}
