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
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	BuildDir   string `json:"build_dir"`
	Icon       Icon   `json:"icon,omitempty"`
}

type Icon struct {
	Linux   string `json:"linux"`
	MacOS   string `json:"darwin"`
	Windows string `json:"windows"`
}

type bundler interface {
	build(v *Version)
}

func Build(os, arch string) error {
	t := time.Now()
	b, e := getBundler(os, arch)
	if e != nil {
		return e
	}
	v := loadVersion()
	b.build(v)
	d := time.Since(t)
	fmt.Printf("build complete in %v seconds\n", d.Seconds())
	return nil
}

func getBundler(os, arch string) (bundler, error) {
	switch os {
	case "darwin":
		return &osx{arch: arch}, nil
	case "windows":
		return &windows{arch: arch}, nil
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
