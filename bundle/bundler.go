package bundle

import (
	"bytes"
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
	e = compile(b, build, platform, name, arch)
	if e != nil {
		return time.Duration(0), e
	}

	return time.Since(start), nil
}

func builder(platform string) (Builder, error) {
	switch platform {
	case "darwin":
		return &darwin{}, nil
	case "windows":
		return &windows{}, nil
	}
	return nil, errors.New("platform not supported")
}

func compile(b Builder, build *Build, platform *Platform, name, arch string) error {
	if arch == "" {
		return errors.New("invalid arch")
	}

	buildId := fmt.Sprint(time.Now().Unix())
	target := fmt.Sprintf("--targets=%s/%s", name, arch)
	flags := flags(platform)

	fatal(exec.Command("xgo", target, flags, "-out", buildId, ".").Run())
	files, e := ioutil.ReadDir(".")
	fatal(e)

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), buildId) {
			return moveFile(f.Name(), b.ExecPath(build, platform, arch))
		}
	}

	return errors.New("executable not found")
}

func flags(p *Platform) string {
	if len(p.Flags) > 0 {
		b := bytes.Buffer{}
		b.WriteString("'")
		for i, s := range p.Flags {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(s)
		}
		b.WriteString("'")

		flags := b.String()
		if len(flags) > 2 {
			return "-ldflags=" + flags
		}
	}
	return ""
}
