package bundler

import (
	"os"
	"path/filepath"
	"github.com/pkg/errors"
	"io/ioutil"
	"encoding/json"
)

type Version struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Icon       string `json:"icon"`
	BuildDir   string `json:"build_dir"`
}

type bundler interface {
	build(v *Version)
}

func Build(os, arch string) error {
	b, e := getBundler(os, arch)
	if e != nil {
		return e
	}
	v := loadVersion()
	b.build(v)
	return nil
}

func getBundler(os, arch string) (bundler, error) {
	switch os {
	case "darwin":
		return &osx{arch:arch}, nil
	case "windows":
		return &windows{arch:arch}, nil
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