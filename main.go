package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dags-/bundler/bundle"
)

func main() {
	start := time.Now()
	b := bundle.LoadBuildFile()

	bundle.Setup(b)

	for name, plat := range b.Platforms {
		for _, arch := range plat.Arch {
			fmt.Printf("building for: %s/%s\n", name, arch)
			t, e := bundle.Bundle(b, plat, name, arch)
			if e != nil {
				log.Println(e)
				continue
			}
			log.Printf("build complete: %s/%s (%.3f seconds)\n", name, arch, t.Seconds())
		}
	}

	log.Printf("build(s) complete in %.3f seconds\n", time.Since(start).Seconds())
}
