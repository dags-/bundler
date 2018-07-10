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

func (d *darwin) Artifact(b *Build, p *Platform, arch string) string {
	name := fmt.Sprintf("%s-%s-%s.app", b.Name, b.Version, arch)
	return filepath.Join(b.Output, "darwin", name)
}

func (d *darwin) ExecPath(b *Build, p *Platform, arch string) string {
	return filepath.Join(d.Artifact(b, p, arch), "Contents", "MacOS", b.Name)
}

func (d *darwin) WriteManifest(b *Build, p *Platform, arch string) error {
	name := fmt.Sprintf("%s-%s-%s.app", b.Name, b.Version, arch)
	path := filepath.Join(b.Output, "darwin", name, "Contents", "Info.plist")
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

func (d *darwin) WriteIcon(b *Build, p *Platform, arch string) error {
	name := fmt.Sprintf("%s-%s-%s.app", b.Name, b.Version, arch)
	_, icon := filepath.Split(p.Icon)
	path := filepath.Join(b.Output, "darwin", name, "Contents", "Resources", icon)
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
