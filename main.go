package main

import (
	"github.com/dags-/bundler/bundler"
	"flag"
	"runtime"
)

func main() {
	os := flag.String("os", runtime.GOOS, "target operating system")
	arch := flag.String("arch", runtime.GOARCH, "target system architecture")
	flag.Parse()

	e := bundler.Build(*os, *arch)
	if e != nil {
		panic(e)
	}
}
