package bundle

import (
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
	args = addArg(args, "-targets", target)
	args = addArg(args, "-out", buildId)
	args = addArg(args, "-ldflags", b.Flags...)
	args = append(args, ".")
	return "xgo", args
}

func addArg(a []string, name string, val ...string) []string {
	value := strings.Join(val, " ")
	if len(val) > 0 {
		value = "'" + value + "'"
	}
	if value != "" {
		a = append(a, name)
		a = append(a, value)
	}
	return a
}
