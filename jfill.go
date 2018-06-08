package jfill

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	scan "github.com/mattn/go-scan"
	"github.com/pkg/errors"
)

const (
	exitOK = iota
	exitError
)

func Run(argv []string) int {
	err := run(argv)
	if err != nil {
		log.Println(err)
		if extr, ok := err.(interface{ ExitCode() int }); ok {
			return extr.ExitCode()
		}
		return exitError
	}
	return exitOK
}

func run(argv []string) error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	var tree interface{}
	if stat.Mode()&os.ModeCharDevice == 0 {
		// input from pipe
		if err := json.NewDecoder(os.Stdin).Decode(&tree); err != nil {
			return err
		}
	}

	var cmdArgs []string
	for _, arg := range argv {
		str, err := fill(arg, tree)
		if err != nil {
			return err
		}
		cmdArgs = append(cmdArgs, str)
	}

	if len(cmdArgs) == 0 {
		return nil
	}
	return runCmd(cmdArgs)
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
