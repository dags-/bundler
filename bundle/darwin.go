package bundle

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
)

type InfoPlist struct {
	Executable string
	Version    string
	Icon       string
	Identifier string
}

type darwin struct {
}

func (d *darwin) Generate() error {
	return exec.Command("go", "generate").Run()
}

func (d *darwin) ExecPath(v *Version, arch string) string {
	name := fmt.Sprintf("%s-%s-%s.App", v.Name, v.Version, arch)
	return filepath.Join(v.BuildDir, "darwin", name, "Content", "MacOS")
}

func (d *darwin) WriteManifest(v *Version, arch string) error {
	name := fmt.Sprintf("%s-%s-%s.app", v.Name, v.Version, arch)
	path := filepath.Join(v.BuildDir, "darwin", name, "Content", "Info.plist")
	mustFile(path)
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	defer f.Close()
	_, icon := filepath.Split(v.Icon.MacOS)
	return template.Must(template.New("info").Parse(infoPlist)).Execute(f, &InfoPlist{
		Executable: v.Name,
		Version:    v.Version,
		Icon:       icon,
		Identifier: v.Identifier,
	})
}

func (d *darwin) WriteIcon(v *Version, arch string) error {
	name := fmt.Sprintf("%s-%s-%s.app", v.Name, v.Version, arch)
	_, icon := filepath.Split(v.Icon.MacOS)
	path := filepath.Join(v.BuildDir, "darwin", name, "Content", "Resources", icon)
	return copyFile(v.Icon.MacOS, path)
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
