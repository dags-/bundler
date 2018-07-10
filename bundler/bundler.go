package bundler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type Version struct {
	Identifier string              `json:"identifier"`
	Name       string              `json:"name"`
	Version    string              `json:"version"`
	BuildDir   string              `json:"build_dir"`
	Targets    map[string][]string `json:"targets"`
	Icon       Icon                `json:"icon,omitempty"`
}

type Icon struct {
	Linux   string `json:"linux"`
	MacOS   string `json:"darwin"`
	Windows string `json:"windows"`
}

type Target struct {
	OS string
}

type bundler interface {
	build(v *Version, arch string)
}

func Build() error {
	start := time.Now()
	v := loadVersion()

	for platform, arch := range v.Targets {
		b, e := getBundler(platform)
		if e != nil {
			log.Println(e)
			continue
		}
		t := time.Now()
		b.build(v)
		d := time.Since(t)
		fmt.Printf("build complete: %s/%s (%.3f seconds)\n", platform, arch, d.Seconds())
	}

	dur := time.Since(start)
	fmt.Printf("total build time: %.3f seconds\n", dur.Seconds())
	return nil
}

func getBundler(platform string) (bundler, error) {
	switch platform {
	case "darwin":
		return &osx{}, nil
	case "windows":
		return &windows{}, nil
	default:
		return nil, errors.New("unsupported os")
	}
}

func loadVersion() *Version {
	var version Version

	d, e := ioutil.ReadFile("version.json")
	fatal(e)

	e = json.Unmarshal(d, &version)
	fatal(e)

	return &version
}

func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

func mustDir(path ...string) string {
	dir := filepath.Join(path...)
	e := os.MkdirAll(dir, os.ModePerm)
	fatal(e)
	return dir
}

func copyFile(from, to string) {
	in, e := os.Open(from)
	if e != nil {
		log.Println(e)
		return
	}
	defer in.Close()

	out, e := os.Create(to)
	if e != nil {
		log.Println(e)
		return
	}
	defer out.Close()

	_, e = io.Copy(out, in)
	if e != nil {
		log.Println(e)
		return
	}
}
