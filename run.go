// +build !windows

package jfill

import (
	"os"
	"os/exec"
	"syscall"
)

func runCmd(cmdArgs []string) error {
	exe, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		return err
	}
	return syscall.Exec(exe, cmdArgs, os.Environ())
}
