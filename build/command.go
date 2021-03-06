package build

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func cmd(cmd string) error {
	log.Println("executing command:", cmd)
	parts := strings.Split(cmd, " ")
	name := parts[0]
	args := parts[1:]
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	return c.Run()
}

func compileCmd(build *Build, buildId, platform, arch string) (cmd string, args []string) {
	if runtime.GOOS == platform {
		return nativeCompile(build, buildId)
	}
	return crossCompile(build, buildId, platform+"/"+arch)
}

func crossCompile(b *Build, buildId, target string) (cmd string, args []string) {
	if b.goVersion != "" {
		args = addArg(args, "-go", b.goVersion)
	}
	args = addArg(args, "-targets", target)
	args = addArg(args, "-out", buildId)
	args = addArg(args, "-ldflags", b.Flags...)
	args = append(args, ".")
	return "xgo", args
}

func nativeCompile(b *Build, buildId string) (cmd string, args []string) {
	args = append(args, "build")
	args = addArg(args, "-o", buildId)
	args = addArg(args, "-ldflags", b.Flags...)
	return "go", args
}

func addArg(a []string, name string, val ...string) []string {
	value := strings.Join(val, " ")
	if len(val) > 1 {
		value = "'" + value + "'"
	}
	if value != "" {
		a = append(a, name)
		a = append(a, value)
	}
	return a
}
