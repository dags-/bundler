package build

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

type BuildScript struct {
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Identifier string            `json:"identifier"`
	Icon       string            `json:"icon"`
	Output     string            `json:"output"`
	Setup      []string          `json:"setup"`
	Targets    map[string]*Build `json:"targets"`
	winIcon    string
	macIcon    string
}

type Build struct {
	Icon     string   `json:"icon"`
	Compress bool     `json:"compress"`
	Generate []string `json:"generate"`
	Flags    []string `json:"flags"`
}

func Setup(script *BuildScript) {
	log.Println("executing setup commands...")
	for _, c := range script.Setup {
		cmd(c)
	}

	if script.Icon != "" {
		log.Println("generating icons...")
		src, e := loadImage(script.Icon)
		if e != nil {
			log.Println("icon error:", e)
			return
		}
		if path, e := writeIcon(convertIco, src, filepath.Join(script.Output, "icon.ico")); e == nil {
			script.winIcon = path
		}
		if path, e := writeIcon(convertIcns, src, filepath.Join(script.Output, "icon.icns")); e == nil {
			script.macIcon = path
		}
	}
}

func LoadBuildFile() *BuildScript {
	var build BuildScript

	d, e := ioutil.ReadFile("build.json")
	fatal(e)

	e = json.Unmarshal(d, &build)
	fatal(e)

	return &build
}
