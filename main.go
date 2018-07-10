package main

import (
	"log"
	"time"

	"github.com/dags-/bundler/bundle"
)

func main() {
	start := time.Now()
	v := bundle.LoadVersion()
	for plat, arch := range v.Targets {
		for _, a := range arch {
			t, e := bundle.Build(plat, a, v)
			if e != nil {
				log.Println(e)
				continue
			}
			log.Printf("build complete: %s/%s (%.3f seconds)\n", plat, a, t.Seconds())
		}
	}
	log.Printf("build(s) complete in %.3f seconds\n", time.Since(start).Seconds())
}
