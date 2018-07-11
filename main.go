package main

import (
	"log"
	"time"

	"github.com/dags-/bundler/bundle"
)

func main() {
	start := time.Now()
	script := bundle.LoadBuildFile()

	log.Println("Setting up")
	bundle.Setup(script)

	log.Println("Running builds")
	for target, build := range script.Targets {
		log.Printf("Target: %s\n", target)
		bundle.Generate(build)
		t, e := bundle.Bundle(script, build, target)
		if e != nil {
			log.Println(" build error:", e)
			continue
		}
		log.Printf(" build complete: %s (%.3f seconds)\n", target, t.Seconds())
	}

	log.Printf("Build(s) complete in %.3f seconds\n", time.Since(start).Seconds())
}
