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
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Println(err)
		return exitError
	}
	if stat.Mode()&os.ModeCharDevice != 0 {
		log.Println("no input from stdin")
		return exitOK
	}

	var tree interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&tree); err != nil {
		log.Println(err)
		return exitError
	}

	var age int
	if err := scan.ScanTree(tree, "/age", &age); err != nil {
		log.Println(err)
		return exitError
	}

	log.Println(age)
	for _, arg := range argv {
		str, err := fill(arg, tree)
		if err != nil {
			log.Println(err)
			return exitError
		}
		log.Println(str)
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
