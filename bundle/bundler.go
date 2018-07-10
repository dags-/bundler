package bundle

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Builder interface {
	ExecPath(build *Build, platform *Platform, arch string) string

	WriteManifest(build *Build, platform *Platform, arch string) error

	WriteIcon(build *Build, platform *Platform, arch string) error
}

func Setup(build *Build) {
	log.Println("executing setup commands...")
	for _, c := range build.Setup {
		cmd(c)
	}
}

func Bundle(build *Build, platform *Platform, name string, arch string) (time.Duration, error) {
	start := time.Now()

	b, e := builder(name)
	if e != nil {
		return time.Duration(0), e
	}

	log.Println("writing manifest...")
	e = b.WriteManifest(build, platform, arch)
	if e != nil {
		return time.Duration(0), e
	}

	log.Println("writing icon...")
	e = b.WriteIcon(build, platform, arch)
	if e != nil {
		log.Println(e)
	}

	log.Println("executing pre build commands...")
	for _, cmd := range platform.Generate {
		exec.Command(cmd).Run()
	}

	log.Println("compiling executable...")
	e = compile(b, build, platform, arch)
	if e != nil {
		return time.Duration(0), e
	}

	return time.Since(start), nil
}

func builder(platform string) (Builder, error) {
	switch platform {
	case "darwin":
		return &darwin{}, nil
	case "linux":

	case "windows":

	}
	return nil, errors.New("platform not supported")
}

func compile(b Builder, build *Build, platform *Platform, arch string) error {
	if arch == "" {
		return errors.New("invalid arch")
	}

	buildId := fmt.Sprint(time.Now().Unix())
	targets := fmt.Sprintf("--targets=%s/%s", platform, arch)

	fatal(exec.Command("xgo", targets, "-out", buildId, ".").Run())

	files, e := ioutil.ReadDir(".")
	fatal(e)

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), buildId) {
			return moveFile(f.Name(), b.ExecPath(build, platform, arch))
		}
	}

	return errors.New("executable not found")
}
