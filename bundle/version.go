package bundle

import (
	"encoding/json"
	"io/ioutil"
)

type Version struct {
	Identifier string              `json:"identifier"`
	Name       string              `json:"name"`
	Version    string              `json:"version"`
	BuildDir   string              `json:"build_dir"`
	Targets    map[string][]string `json:"targets"`
	Icon       Icon                `json:"icon,omitempty"`
}

type Icon struct {
	Linux   string `json:"linux"`
	MacOS   string `json:"darwin"`
	Windows string `json:"windows"`
}

func LoadVersion() *Version {
	var version Version

	d, e := ioutil.ReadFile("version.json")
	fatal(e)

	e = json.Unmarshal(d, &version)
	fatal(e)

	return &version
}
