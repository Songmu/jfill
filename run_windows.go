package jfill

import (
	"os"
	"os/exec"

	"github.com/Songmu/wrapcommander"
)

func runCmd(argv []string) error {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Wait(); err != nil {
		os.Exit(wrapcommander.ResolveExitCode(err))
	}
	return nil
}
