//+build !linux

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
	//c.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGHUP}
	return c
}
