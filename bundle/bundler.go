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
	Generate() error

	ExecPath(*Version, string) string

	WriteManifest(*Version, string) error

	WriteIcon(*Version, string) error
}

func Build(platform, arch string, v *Version) (time.Duration, error) {
	start := time.Now()

	b, e := builder(platform)
	if e != nil {
		return time.Duration(0), e
	}

	// need to do this before go:generate for windows platform
	log.Println("writing manifest...")
	e = b.WriteManifest(v, arch)
	if e != nil {
		return time.Duration(0), e
	}

	log.Println("writing icon...")
	e = b.WriteIcon(v, arch)
	if e != nil {
		log.Println(e)
	}

	log.Println("performing go:generate")
	e = b.Generate()
	if e != nil {
		return time.Duration(0), e
	}

	log.Println("building executable")
	e = build(b, v, platform, arch)
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

func build(b Builder, v *Version, platform, arch string) error {
	if platform == "" || arch == "" {
		return errors.New("invalid platform/arch")
	}

	buildId := fmt.Sprint(time.Now().Unix())
	targets := fmt.Sprintf("--targets=%s/%s", platform, arch)

	fatal(exec.Command("xgo", targets, "-out", buildId, ".").Run())

	files, e := ioutil.ReadDir(".")
	fatal(e)

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), buildId) {
			return moveFile(f.Name(), b.ExecPath(v, arch))
		}
	}

	return errors.New("executable not found")
}
