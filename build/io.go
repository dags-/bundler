package build

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var wd = flag.String("workdir", ".", "Set the workdir")

func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

func WorkDir() string {
	if *wd == "." {
		return "."
	}

	path := *wd
	defer log.Println("found work dir:", path)

	if exists(path) {
		return path
	}

	path = filepath.Join(os.Getenv("GOPATH"), "src", path)
	if exists(path) {
		return path
	}

	log.Println("go getting project:", path)
	e := exec.Command("go", "get", "-u", path).Run()
	if e == nil && exists(path) {
		return path
	}

	panic("invalid target directory: " + path)
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

func download(url, path string) (*os.File, error) {
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	mustFile(path)
	o, e := os.Create(path)
	if e != nil {
		return nil, e
	}

	_, e = io.Copy(o, r.Body)
	return o, e
}

func applyTemplate(text, path string, i interface{}) error {
	mustFile(path)
	o, e := os.Create(path)
	if e != nil {
		return e
	}
	defer o.Close()
	t := template.Must(template.New("template").Parse(text))
	return t.Execute(o, i)
}
