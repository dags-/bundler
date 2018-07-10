package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(path, platform string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	dir, name := filepath.Split(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = fmt.Sprintf("%s-%s.zip", platform, name)
	zipPath := filepath.Join(dir, name)

	out, e := os.Create(zipPath)
	if e != nil {
		panic(e)
	}
	defer out.Close()

	wr := zip.NewWriter(out)
	defer wr.Close()
	if file.IsDir() {
		return filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
			return writeToZip(wr, path, path, file)
		})
	} else {
		return writeToZip(wr, path, file.Name(), file)
	}
}

func writeToZip(wr *zip.Writer, from string, to string, file os.FileInfo) error {
	if file.IsDir() {
		return nil
	}

	in, e := os.Open(from)
	if e != nil {
		return e
	}
	defer in.Close()

	w, e := wr.Create(to)
	if e != nil {
		return e
	}

	_, e = io.Copy(w, in)
	return e
}
