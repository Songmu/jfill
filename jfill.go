package jfill

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	scan "github.com/mattn/go-scan"
	"github.com/pkg/errors"
)

const (
	exitOK = iota
	exitError
)

func Run(argv []string) int {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Println(err)
		return exitError
	}

	var tree interface{}
	if stat.Mode()&os.ModeCharDevice == 0 {
		// input from pipe
		if err := json.NewDecoder(os.Stdin).Decode(&tree); err != nil {
			log.Println(err)
			return exitError
		}
	}

	var cmdArgs []string
	for _, arg := range argv {
		str, err := fill(arg, tree)
		if err != nil {
			log.Println(err)
			return exitError
		}
		cmdArgs = append(cmdArgs, str)
	}

	if len(cmdArgs) == 0 {
		return exitOK
	}
	exe, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		log.Println(err)
		return exitError
	}
	if err := syscall.Exec(exe, cmdArgs, os.Environ()); err != nil {
		log.Println(err)
		return exitError
	}
	return exitOK
}

var fillReg = regexp.MustCompile(`\{\{` +
	`([a-zA-Z0-9/\[\]]+)` +
	`(?::([^{}]*))?` +
	`\}\}`)

func fill(str string, tree interface{}) (string, error) {
	var retErr error
	ret := fillReg.ReplaceAllStringFunc(str, func(match string) string {
		m := fillReg.FindStringSubmatch(match)
		if len(m) < 3 {
			retErr = fmt.Errorf("something went wrong")
			return ""
		}
		jpath := m[1]
		if !strings.HasPrefix(jpath, "/") {
			jpath = "/" + jpath
		}
		var in interface{}
		if err := scan.ScanTree(tree, jpath, &in); err != nil {
			if m[2] == "" {
				retErr = errors.Wrap(err, "failed to scan")
				return ""
			}
			return m[2]
		}
		return fmt.Sprintf("%v", in)
	})
	return ret, retErr
}
