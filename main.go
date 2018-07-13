package main

import (
	"log"
	"os"
	"time"

	"github.com/dags-/bundler/build"
)

func main() {
	log.SetPrefix("[build] ")
	log.Println("SETTING WORK DIR")
	os.Chdir(build.WorkDir())

	start := time.Now()
	script := build.LoadBuildFile()

	log.Println("CLEANING BUILD DIR")
	os.RemoveAll(script.Output)

	log.Println("RUNNING SETUP SCRIPTS")
	build.Setup(script)

	log.Println("RUNNING BUILDS")
	for target, b := range script.Targets {
		log.SetPrefix("[" + target + "]")
		log.Printf("building for: %s\n", target)
		t, e := build.Run(script, b, target)
		if e != nil {
			log.Println("build error:", e)
			continue
		}
		log.Printf("build complete: %.3f seconds\n", t.Seconds())
	}

	log.SetPrefix("[build] ")
	log.Printf("BUILD(S) COMPLETE: %.3f SECONDS\n", time.Since(start).Seconds())
}
