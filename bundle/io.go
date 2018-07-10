package bundle

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func cmd(cmd string) error {
	parts := strings.Split(cmd, " ")
	name := parts[0]
	args := parts[1:]
	return exec.Command(name, args...).Run()
}

func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

func exists(path string) bool {
	_, e := os.Stat(path)
	return e == nil
}

func mustFile(path string) error {
	return mustDir(filepath.Dir(path))
}

func mustDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func moveFile(from, to string) error {
	if !exists(from) {
		return os.ErrNotExist
	}
	e := mustFile(to)
	if e != nil {
		return e
	}
	return os.Rename(from, to)
}

func copyFile(from, to string) error {
	in, e := os.Open(from)
	if e != nil {
		return e
	}
	defer in.Close()

	e = mustFile(to)
	if e != nil {
		return e
	}

	out, e := os.Create(to)
	if e != nil {
		return e
	}
	defer out.Close()

	_, e = io.Copy(out, in)
	if e != nil {
		return e
	}

	return nil
}
