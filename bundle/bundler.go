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
	Artifact(build *BuildScript, platform *Build, arch string) string

	ExecPath(build *BuildScript, platform *Build, arch string) string

	WriteManifest(build *BuildScript, platform *Build, arch string) error

	WriteIcon(build *BuildScript, platform *Build, arch string) error
}

func Setup(build *BuildScript) {
	log.Println(" executing setup commands...")
	for _, c := range build.Setup {
		cmd(c)
	}
}

func Generate(build *Build) {
	log.Println(" executing generate commands...")
	for _, c := range build.Generate {
		cmd(c)
	}
}

func Bundle(script *BuildScript, build *Build, target string) (time.Duration, error) {
	platform, arch := splitTarget(target)
	if platform == "" || arch == "" {
		return time.Duration(0), errors.New("invalid target " + target)
	}

	b, e := builder(platform)
	if e != nil {
		return time.Duration(0), e
	}

	start := time.Now()
	log.Println(" writing manifest...")
	e = b.WriteManifest(script, build, arch)
	if e != nil {
		return time.Duration(0), e
	}

	log.Println(" writing icon...")
	e = b.WriteIcon(script, build, arch)
	if e != nil {
		log.Println(e)
	}

	log.Println(" compiling executable...")
	e = compile(b, script, build, platform, arch)
	if e != nil {
		return time.Duration(0), e
	}

	if build.Compress {
		log.Println(" compressing executable")
		e = Compress(b.Artifact(script, build, arch), platform)
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

func compile(b Builder, script *BuildScript, build *Build, platform, arch string) error {
	if arch == "" {
		return errors.New("invalid arch")
	}

	buildId := fmt.Sprint(time.Now().Unix())
	targets := fmt.Sprintf("--targets=%s/%s", platform, arch)
	flags := flags(build)
	args := []string{targets, "-out", buildId}
	if len(flags) > 0 {
		args = append(args, flags)
	}
	args = append(args, ".")

	// todo: remove
	log.Println("  (debug) compile command: xgo", strings.Join(args, " "))

	fatal(exec.Command("xgo", args...).Run())
	files, e := ioutil.ReadDir(".")
	fatal(e)

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), buildId) {
			return moveFile(f.Name(), b.ExecPath(script, build, toNormal(arch)))
		}
	}

	return errors.New("executable not found")
}

func splitTarget(s string) (string, string) {
	parts := strings.Split(s, "/")
	if len(parts) == 1 {
		return parts[0], "*"
	}
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return s, ""
}

func flags(b *Build) string {
	if len(b.Flags) > 0 {
		buf := bytes.Buffer{}
		for i, s := range b.Flags {
			if i > 0 {
				buf.WriteString(" ")
			}
			buf.WriteString(s)
		}

		flags := buf.String()
		if len(flags) > 0 {
			return "-ldflags='" + flags + "'"
		}
	}
	return ""
}
