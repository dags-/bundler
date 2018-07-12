package build

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type builder interface {
	init(script *BuildScript, build *Build, arch string)

	preCompile() error

	postCompile() error

	artifact() string

	executable() string
}

func Run(script *BuildScript, build *Build, target string) (time.Duration, error) {
	parts := strings.Split(target, "/")
	if len(parts) != 2 {
		return time.Duration(0), errors.New("invalid target")
	}

	platform, arch := parts[0], parts[1]
	b, e := getBuilder(platform)
	if e != nil {
		return time.Duration(0), e
	}

	start := time.Now()
	if build.Icon == "" {
		build.Icon = icon(script, platform)
	}

	log.Println("INIT")
	b.init(script, build, toNormal(arch))

	log.Println("PRE-COMPILE")
	if e := b.preCompile(); e != nil {
		return time.Duration(0), e
	}

	log.Println("GENERATE")
	for _, c := range build.Generate {
		cmd(c)
	}

	log.Println("COMPILE")
	if e := compile(b, build, platform, arch); e != nil {
		return time.Duration(0), e
	}

	log.Println("POST-COMPILE")
	if e := b.postCompile(); e != nil {
		return time.Duration(0), e
	}

	if build.Compress {
		log.Println("ARCHIVE")
		if e := compress(b.artifact(), platform); e != nil {
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
	buildId := fmt.Sprint(time.Now().Unix())
	cmd, args := compileCmd(build, buildId, platform, arch)
	log.Printf("compile command: %s %s\n", cmd, strings.Join(args, " "))

	if e := exec.Command(cmd, args...).Run(); e != nil {
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
