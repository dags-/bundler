package build

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func compress(path, platform string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	platform = toNormal(platform)
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
		return filepath.Walk(path, func(from string, file os.FileInfo, err error) error {
			to, err := filepath.Rel(dir, from)
			if err != nil {
				log.Println("compress error:", err)
				return nil
			}
			return writeToZip(wr, from, to, file)
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
