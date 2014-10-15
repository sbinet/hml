//+build windows

package main

import (
	"os"
	"os/exec"
)

func (code Code) command(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Dir = code.Root
	//c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return c
}

func (code Code) kill(p *os.Process) error {
	//FIXME(sbinet): on windows, there is AFAIK no way to kill a whole hierarchy of processes
	// because no parent/child(ren) concept
	return p.Kill()
}
