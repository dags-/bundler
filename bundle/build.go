package bundle

import (
	"encoding/json"
	"io/ioutil"
)

type Build struct {
	Name       string               `json:"name"`
	Version    string               `json:"version"`
	Identifier string               `json:"identifier"`
	Output     string               `json:"output"`
	Setup      []string             `json:"setup"`
	Platforms  map[string]*Platform `json:"platforms"`
}

type Platform struct {
	Icon     string   `json:"icon"`
	Arch     []string `json:"arch"`
	Generate []string `json:"pre"`
	Flags    []string `json:"flags"`
	Compress bool     `json:"compress"`
}

func LoadBuildFile() *Build {
	var build Build

	d, e := ioutil.ReadFile("build.json")
	fatal(e)

	e = json.Unmarshal(d, &build)
	fatal(e)

	return &build
}
