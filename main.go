package main

import (
	"log"
	"time"

	"github.com/dags-/bundler/bundle"
)

func main() {
	start := time.Now()
	b := bundle.LoadBuildFile()

	log.Println("setting up")
	bundle.Setup(b)

	log.Println("running build tasks")
	for name, plat := range b.Platforms {
		bundle.Generate(plat)
		for _, arch := range plat.Arch {
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
