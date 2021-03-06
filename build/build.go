package build

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type builder interface {
	init(script *Script, build *Build, arch string)

	artifact() string

	executable() string

	preCompile() error

	postCompile() error

	compress() error

	clean()
}

func Native(target string) bool {
	parts := strings.Split(target, "/")
	if len(parts) != 2 {
		return false
	}
	return parts[0] == runtime.GOOS
}

func Run(script *Script, build *Build, target string) (time.Duration, error) {
	parts := strings.Split(target, "/")
	if len(parts) != 2 {
		return time.Duration(0), errors.New("invalid target")
	}

	platform, arch := parts[0], parts[1]
	b, e := getBuilder(platform)
	if e != nil {
		return time.Duration(0), e
	}
	defer b.clean()

	start := time.Now()
	if build.Icon == "" {
		build.Icon = icon(script, platform)
	}

	log.Println("# INIT")
	b.init(script, build, toNormal(arch))

	log.Println("# PRE-COMPILE")
	if e := b.preCompile(); e != nil {
		return time.Duration(0), e
	}

	log.Println("# GENERATE")
	for _, c := range build.Generate {
		cmd(c)
	}

	log.Println("# COMPILE")
	build.goVersion = script.GoVersion
	if e := compile(b, build, platform, arch); e != nil {
		return time.Duration(0), e
	}

	log.Println("# POST-COMPILE")
	if e := b.postCompile(); e != nil {
		return time.Duration(0), e
	}

	if build.Compress {
		log.Println("# COMPRESS")
		if e := b.compress(); e != nil {
			return time.Duration(0), e
		}
	}

	return time.Since(start), nil
}

func getBuilder(platform string) (builder, error) {
	switch platform {
	case "linux":
		return &linux{}, nil
	case "darwin":
		return &darwin{}, nil
	case "windows":
		return &windows{}, nil
	default:
		return nil, errors.New("unsupported platform")
	}
}

func compile(b builder, build *Build, platform, arch string) error {
	e := os.Setenv("GOOS", platform)
	if e != nil {
		log.Println("set env error:", e)
	}

	e = os.Setenv("GOARCH", arch)
	if e != nil {
		log.Println("set env error:", e)
	}

	buildId := fmt.Sprint(time.Now().Unix())
	cmd, args := compileCmd(build, buildId, platform, arch)
	log.Printf("compile command: %s %s\n", cmd, strings.Join(args, " "))
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if e := c.Run(); e != nil {
		return e
	}

	files, e := ioutil.ReadDir(".")
	if e != nil {
		return e
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), buildId) {
			return moveFile(f.Name(), b.executable())
		}
	}

	return errors.New("executable not found")
}

func toNormal(name string) string {
	switch name {
	case "darwin":
		return "macOS"
	case "amd64":
		return "x64"
	case "386":
		return "x32"
	}
	return name
}
