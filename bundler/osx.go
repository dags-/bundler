package bundler

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Manifest struct {
	Executable string
	Version    string
	Icon       string
	Identifier string
}

type osx struct {
	arch string
}

func (o *osx) build(version *Version) {
	buildDir := mustDir(version.BuildDir, "darwin")
	contents := mustDir(buildDir, version.Name+".app", "Contents")

	log.Println("running go generate")
	o.generate()

	log.Println("running go build")
	o.executable(contents, version.Name)

	log.Println("writing Info.plist")
	o.manifest(contents, version)

	log.Println("copying icon file")
	o.icon(contents, version.Icon.MacOS)
}

func (o *osx) generate() {
	fatal(exec.Command("go", "generate").Run())
}

func (o *osx) executable(contentsDir, name string) {
	execFile := filepath.Join(mustDir(contentsDir, "MacOS"), name)
	fatal(exec.Command("go", "build", "-o", execFile).Run())
}

func (o *osx) manifest(contentsDir string, version *Version) {
	infoFile := filepath.Join(contentsDir, "Info.plist")
	info, e := os.Create(infoFile)
	fatal(e)
	defer info.Close()

	_, iconName := filepath.Split(version.Icon.MacOS)
	t := template.Must(template.New("info").Parse(templ))
	fatal(t.Execute(info, &Manifest{
		Executable: version.Name,
		Version:    version.Version,
		Icon:       iconName,
		Identifier: version.Identifier,
	}))
}

func (o *osx) icon(contentsDir, icon string) {
	if icon == "" {
		return
	}

	_, iconName := filepath.Split(icon)
	path := filepath.Join(mustDir(contentsDir, "Resources"), iconName)
	copyFile(icon, path)
}

const templ = `<?xml version="1.0" encoding="UTF-8"?>
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
