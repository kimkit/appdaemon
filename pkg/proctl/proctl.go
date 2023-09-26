package proctl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var (
	ErrProcessExitWaitingTimeout = fmt.Errorf("process exit waiting timeout")
)

type Process struct {
	bin          string
	args         []string
	env          map[string]string
	output       string
	outputWriter io.Writer
	cmd          *exec.Cmd
	status       int64
	mu           sync.Mutex
	fp           *os.File
}

func NewProcess(bin string, args ...string) *Process {
	return &Process{
		bin:  bin,
		args: args,
		env:  make(map[string]string),
	}
}

func (p *Process) SetEnv(name, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.env[name] = value
}

func (p *Process) SetOutput(file string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.output = file
}

func (p *Process) SetOutputWriter(writer io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.outputWriter = writer
}

func (p *Process) wait() {
	if err := p.cmd.Wait(); err != nil {
		// pass
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.status == 1 {
		p.status = 0
		if p.fp != nil {
			if err := p.fp.Close(); err != nil {
				// pass
			}
		}
	}
}

func (p *Process) Run() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.status == 0 {
		p.cmd = exec.Command(p.bin, p.args...)
		if p.outputWriter != nil {
			p.cmd.Stdout = p.outputWriter
			p.cmd.Stderr = p.outputWriter
		} else {
			p.fp = nil
			if p.output != "" {
				fp, err := os.OpenFile(p.output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					// pass
				} else {
					p.fp = fp
				}
			}
			p.cmd.Stdout = p.fp
			p.cmd.Stderr = p.fp
		}
		var env []string
		for _, v := range os.Environ() {
			env = append(env, v)
		}
		for k, v := range p.env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		p.cmd.Env = env
		p.cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
		if err := p.cmd.Start(); err != nil {
			return err
		}
		go p.wait()
		p.status = 1
	}
	return nil
}

func (p *Process) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.status == 1
}

func (p *Process) Kill() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.status == 1 {
		return p.cmd.Process.Kill()
	}
	return nil
}

func (p *Process) Signal(sig syscall.Signal) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.status == 1 {
		return p.cmd.Process.Signal(sig)
	}
	return nil
}

func (p *Process) Wait(second int) error {
	n := 0
	for {
		if p.IsRunning() {
			time.Sleep(time.Millisecond * 100)
			n++
			if second > 0 && n > second*10 {
				return ErrProcessExitWaitingTimeout
			}
		} else {
			return nil
		}
	}
}
