package bundle

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type InfoPlist struct {
	Executable string
	Version    string
	Icon       string
	Identifier string
}

type darwin struct{}

func (d *darwin) Artifact(b *BuildScript, p *Build, arch string) string {
	name := fmt.Sprintf("%s-%s-%s.app", b.Name, b.Version, toNormal(arch))
	return filepath.Join(b.Output, "darwin", name)
}

func (d *darwin) ExecPath(b *BuildScript, p *Build, arch string) string {
	return filepath.Join(d.Artifact(b, p, toNormal(arch)), "Content", "MacOS", b.Name)
}

func (d *darwin) WriteManifest(b *BuildScript, p *Build, arch string) error {
	name := fmt.Sprintf("%s-%s-%s.app", b.Name, b.Version, toNormal(arch))
	path := filepath.Join(b.Output, "darwin", name, "Content", "Info.plist")
	mustFile(path)
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	defer f.Close()
	_, icon := filepath.Split(p.Icon)
	return template.Must(template.New("info").Parse(infoPlist)).Execute(f, &InfoPlist{
		Executable: b.Name,
		Version:    b.Version,
		Icon:       icon,
		Identifier: b.Identifier,
	})
}

func (d *darwin) WriteIcon(b *BuildScript, p *Build, arch string) error {
	name := fmt.Sprintf("%s-%s-%s.app", b.Name, b.Version, toNormal(arch))
	_, icon := filepath.Split(p.Icon)
	path := filepath.Join(b.Output, "darwin", name, "Content", "Resources", icon)
	return copyFile(p.Icon, path)
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
