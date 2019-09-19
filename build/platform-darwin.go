package build

import (
	"fmt"
	"log"
	"path/filepath"
)

type darwin struct {
	*Build
	*Script
	appDir   string
	exePath  string
	infoPath string
	iconPath string
}

type InfoPlist struct {
	Executable string
	Version    string
	Icon       string
	Identifier string
}

func (d *darwin) artifact() string {
	return d.appDir
}

func (d *darwin) executable() string {
	return d.exePath
}

func (d *darwin) init(script *Script, build *Build, arch string) {
	version := fmt.Sprintf("%s-%s", script.Version, arch)
	d.Build = build
	d.Script = script
	d.appDir = filepath.Join(script.Output, "darwin", version, script.Name+".app")
	d.exePath = filepath.Join(d.appDir, "Contents", "MacOS", script.Name)
	d.infoPath = filepath.Join(d.appDir, "Contents", "Info.plist")
	d.iconPath = filepath.Join(d.appDir, "Contents", "Resources", "icon.icns")
}

func (d *darwin) preCompile() error {
	mustFile(d.iconPath)
	mustFile(d.infoPath)

	log.Println("copying icon")
	if e := copyFile(d.Build.Icon, d.iconPath); e != nil {
		log.Println("icon error:", e)
	}

	log.Println("writing info.plist")
	if e := applyTemplate(infoPlist, d.infoPath, d.manifest()); e != nil {
		return e
	}
	return nil
}

func (d *darwin) postCompile() error {
	return nil
}

func (d *darwin) compress() error {
	return compress(d.artifact(), "darwin")
}

func (d *darwin) clean() {}

func (d *darwin) manifest() interface{} {
	return &InfoPlist{
		Executable: d.Name,
		Version:    d.Version,
		Icon:       "icon.icns",
		Identifier: d.Identifier,
	}
}

const infoPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>{{ .Executable }}</string>
	<key>CFBundleIconFile</key>
	<string>{{ .Icon }}</string>
	<key>CFBundleIdentifier</key>
	<string>{{ .Identifier }}</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>LSUIElement</key>
	<true/>
</dict>
</plist>`
