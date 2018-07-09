package bundler

import (
	"text/template"
	"os"
	"path/filepath"
	"os/exec"
	"io"
	"log"
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
	o.icon(contents, version.Icon)
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

	_, iconName := filepath.Split(version.Icon)
	t := template.Must(template.New("info").Parse(templ))
	fatal(t.Execute(info, &Manifest{
		Executable: version.Name,
		Version: version.Version,
		Icon: iconName,
		Identifier: version.Identifier,
	}))
}

func (o *osx) icon(contentsDir, icon string) {
	from, e := os.Open(icon)
	if e != nil {
		log.Println(e)
		return
	}
	defer from.Close()

	_, iconName := filepath.Split(icon)
	path := filepath.Join(mustDir(contentsDir, "Resources"), iconName)
	to, e := os.Create(path)
	if e != nil {
		log.Println(e)
		return
	}
	defer to.Close()

	_, e = io.Copy(to, from)
	if e != nil {
		log.Println(e)
		return
	}
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