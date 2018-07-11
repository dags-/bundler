package bundle

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func cmd(cmd string) error {
	parts := strings.Split(cmd, " ")
	name := parts[0]
	args := parts[1:]
	return exec.Command(name, args...).Run()
}

func compileCmd(b *Build, buildId, target string) (cmd string, args []string) {
	cmd = "xgo"
	args = addArg(args, "-targets", target)
	args = addArg(args, "-out", buildId)
	args = addArg(args, "-ldflags", b.Flags...)
	args = append(args, ".")
	return cmd, args
}

func addArg(a []string, name string, val ...string) []string {
	if len(val) == 0 {
		return a
	}
	var value string
	if len(val) == 1 {
		value = val[0]
	} else {
		buf := bytes.Buffer{}
		for i, s := range val {
			if i > 0 {
				buf.WriteRune(' ')
			}
			buf.WriteString(s)
		}
		value = buf.String()
	}
	if value == "" {
		return a
	}
	a = append(a, name)
	a = append(a, fmt.Sprint("'", value, "'"))
	return a
}
