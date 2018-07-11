package bundle

import (
	"encoding/json"
	"io/ioutil"
)

type BuildScript struct {
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Identifier string            `json:"identifier"`
	Output     string            `json:"output"`
	Setup      []string          `json:"setup"`
	Targets    map[string]*Build `json:"targets"`
}

type Build struct {
	Icon     string   `json:"icon"`
	Compress bool     `json:"compress"`
	Generate []string `json:"generate"`
	Flags    []string `json:"flags"`
}

func LoadBuildFile() *BuildScript {
	var build BuildScript

	d, e := ioutil.ReadFile("build.json")
	fatal(e)

	e = json.Unmarshal(d, &build)
	fatal(e)

	return &build
}
