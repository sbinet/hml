package main

import (
	"os"
	"os/exec"
	"syscall"
)

func (code Code) command(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Dir = code.Root
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return c
}

func (code Code) kill(p *os.Process) error {
	pgid, err := syscall.Getpgid(p.Pid)
	if err != nil {
		return err
	}
	err = syscall.Kill(-pgid, syscall.SIGKILL) // note the minus sign
	return err
}
