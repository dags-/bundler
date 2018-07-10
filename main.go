package main

import (
	"github.com/dags-/bundler/bundler"
)

func main() {
	e := bundler.Build()
	if e != nil {
		panic(e)
	}
}
