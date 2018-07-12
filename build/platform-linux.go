package build

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type linux struct {
	*Build
	*Script
	appDirPath string
	appImgPath string
	exePath    string
	maniPath   string
	iconPath   string
}

type Desktop struct {
	Name       string
	Executable string
	Icon       string
	Categories string
}

func (l *linux) artifact() string {
	return l.appImgPath
}

func (l *linux) executable() string {
	return l.exePath
}

func (l *linux) init(script *Script, build *Build, arch string) {
	name := fmt.Sprintf("%s-%s-%s", script.Name, script.Version, arch)
	l.Build = build
	l.Script = script
	l.appDirPath = filepath.Join(script.Output, "linux", name+".AppDir")
	l.appImgPath = filepath.Join(script.Output, "linux", name+".AppImage")
	l.exePath = filepath.Join(l.appDirPath, "AppRun")
	l.iconPath = filepath.Join(l.appDirPath, script.Name+".png")
	l.maniPath = filepath.Join(l.appDirPath, script.Name+".desktop")
}

func (l *linux) preCompile() error {
	mustFile(l.iconPath)
	mustFile(l.maniPath)

	log.Println("writing icon...")
	if e := copyFile(l.Build.Icon, l.iconPath); e != nil {
		log.Println(" icon error:", e)
	}

	log.Println("writing .desktop...")
	if e := applyTempl(desktop, l.maniPath, l.manifest()); e != nil {
		return e
	}

	return nil
}

func (l *linux) postCompile() error {
	log.Println("packaging appimage...")
	tool, e := l.getImageTool()
	if e != nil {
		return e
	}
	if e := exec.Command(tool, l.appDirPath).Run(); e != nil {
		return e
	}
	return nil
}

func (l *linux) getImageTool() (string, error) {
	path := filepath.Join(l.Output, "linux", "appimagetool-x86_64.AppImage")
	if exists(path) {
		return path, nil
	}

	r, e := http.Get("https://github.com/AppImage/AppImageKit/releases/download/10/appimagetool-x86_64.AppImage")
	if e != nil {
		return "", e
	}
	defer r.Body.Close()

	mustFile(path)
	f, e := os.Create(path)
	if e != nil {
		return "", e
	}
	defer f.Close()

	if _, e = io.Copy(f, r.Body); e != nil {
		return "", e
	}

	f.Chmod(os.ModePerm)

	return path, nil
}

func (l *linux) manifest() interface{} {
	return &Desktop{
		Name:       l.Name,
		Icon:       l.Name,
		Executable: "AppRun",
		Categories: l.MetaData["categories"],
	}
}

const desktop = `[Desktop Entry]
Name={{ .Name }}
Exec={{ .Executable }}
Icon={{ .Icon }}
Type=Application
Categories={{ .Categories }};`
