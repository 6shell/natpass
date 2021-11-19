// +build !windows

package shell

import (
	"errors"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// Exec execute shell command
func (link *Link) Exec() error {
	var cmd *exec.Cmd
	if len(link.parent.cfg.Exec) > 0 {
		cmd = exec.Command(link.parent.cfg.Exec)
	}
	if cmd == nil {
		dir, err := exec.LookPath("bash")
		if err == nil {
			cmd = exec.Command(dir)
		}
	}
	if cmd == nil {
		dir, err := exec.LookPath("sh")
		if err == nil {
			cmd = exec.Command(dir)
		}
	}
	if cmd == nil {
		return errors.New("no shell command supported")
	}
	cmd.Env = append(os.Environ(), link.parent.cfg.Env...)
	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	go cmd.Wait() // defunct process
	link.stdin = f
	link.stdout = f
	link.pid = cmd.Process.Pid
	return nil
}

func (link *Link) onClose() {
	if link.stdin != nil {
		link.stdin.Close()
	}
}

func (link *Link) resize(rows, cols uint32) {
	if link.stdin != nil {
		pty.Setsize(link.stdin.(*os.File), &pty.Winsize{
			Rows: uint16(rows),
			Cols: uint16(cols),
		})
	}
}
