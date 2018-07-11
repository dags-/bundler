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

	log.Println(" executing generate commands...")
	for _, c := range build.Generate {
		cmd(c)
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
	if platform == "" || arch == "" {
		return errors.New("invalid target")
	}

	buildId := fmt.Sprint(time.Now().Unix())
	target := fmt.Sprint(platform, "/", arch)
	cmd, args := compileCmd(build, buildId, target)

	// todo: remove
	log.Printf(" (debug) command: %s %s\n", cmd, strings.Join(args, " "))

	e := exec.Command(cmd, args...).Run()
	fatal(e)

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
