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
	Artifact(build *Build, platform *Platform, arch string) string

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
	name = toInternal(name)
	arch = toInternal(arch)

	b, e := builder(name)
	if e != nil {
		return time.Duration(0), e
	}

	start := time.Now()
	log.Println("writing manifest...")
	e = b.WriteManifest(build, platform, toNormal(arch))
	if e != nil {
		return time.Duration(0), e
	}

	log.Println("writing icon...")
	e = b.WriteIcon(build, platform, toNormal(arch))
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

	if platform.Compress {
		log.Println("compressing executable")
		e = compress(b.Artifact(build, platform, toNormal(arch)), toNormal(name))
	}

	return time.Since(start), e
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
	target := fmt.Sprintf("--targets=%s/%s", toInternal(name), toInternal(arch))
	flags := flags(platform)
	args := []string{target, "-out", buildId}
	if len(flags) > 0 {
		args = append(args, flags)
	}
	args = append(args, ".")

	fatal(exec.Command("xgo", args...).Run())
	files, e := ioutil.ReadDir(".")
	fatal(e)

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), buildId) {
			return moveFile(f.Name(), b.ExecPath(build, platform, toNormal(arch)))
		}
	}

	return errors.New("executable not found")
}

func flags(p *Platform) string {
	if len(p.Flags) > 0 {
		b := bytes.Buffer{}
		for i, s := range p.Flags {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(s)
		}

		flags := b.String()
		if len(flags) > 0 {
			return "-ldflags='" + flags + "'"
		}
	}
	return ""
}
