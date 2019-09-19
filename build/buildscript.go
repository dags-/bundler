package build

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Script struct {
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Identifier string            `json:"identifier"`
	Icon       string            `json:"icon"`
	Output     string            `json:"output"`
	Setup      []string          `json:"setup"`
	GoVersion  string            `json:"go_version"`
	Targets    map[string]*Build `json:"targets"`
	winIcon    string
	macIcon    string
}

type Build struct {
	Icon      string            `json:"icon"`
	Compress  bool              `json:"compress"`
	Generate  []string          `json:"generateInfo"`
	Flags     []string          `json:"flags"`
	MetaData  map[string]string `json:"meta"`
	goVersion string
}

func Setup(script *Script) {
	log.Println("executing setup commands")
	for _, c := range script.Setup {
		cmd(c)
	}

	if script.Icon != "" {
		log.Println("generating icons")
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

	log.Println("generating metainfo")
	e := generateInfo(script.Name, script.Version)
	if e != nil {
		log.Println("failed to generate metainfo/metainfo.go", e)
	}
}

func LoadBuildFile() *Script {
	var build Script

	d, e := ioutil.ReadFile("build.json")
	fatal(e)

	e = json.Unmarshal(d, &build)
	fatal(e)

	return &build
}
