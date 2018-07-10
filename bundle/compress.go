package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func compress(path, platform string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	name := strings.TrimSuffix(path, filepath.Ext(path))
	zipPath := fmt.Sprintf("%s-%s.zip", platform, name)
	out, e := os.Create(zipPath)
	if e != nil {
		panic(e)
	}
	defer out.Close()

	wr := zip.NewWriter(out)
	defer wr.Close()
	if file.IsDir() {
		return filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
			return writeToZip(wr, file)
		})
	} else {
		return writeToZip(wr, file)
	}
}

func writeToZip(wr *zip.Writer, file os.FileInfo) error {
	if file.IsDir() {
		return nil
	}

	w, e := wr.Create(file.Name())
	if e != nil {
		return e
	}

	in, e := os.Open(file.Name())
	if e != nil {
		return e
	}
	defer in.Close()

	_, e = io.Copy(w, in)
	return e
}
