package main

import (
	"log"
	"os"
	"time"

	"github.com/dags-/bundler/build"
)

func main() {
	log.SetPrefix("[build] ")

	start := time.Now()
	script := build.LoadBuildFile()

	log.Println("Cleaning")
	os.RemoveAll(script.Output)

	log.Println("Setting up")
	build.Setup(script)

	log.Println("Running builds")
	for target, b := range script.Targets {
		log.SetPrefix("[" + target + "]")
		log.Printf("building for: %s\n", target)
		t, e := build.Run(script, b, target)
		if e != nil {
			log.Println("build error:", e)
			continue
		}
		log.Printf("build complete: %s (%.3f seconds)\n", target, t.Seconds())
	}

	log.SetPrefix("[build] ")
	log.Printf("Build(s) complete in %.3f seconds\n", time.Since(start).Seconds())
}
