package build

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

type linux struct {
	*Build
	*BuildScript
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

func (l *linux) init(script *BuildScript, build *Build, arch string) {
	name := fmt.Sprintf("%s-%s-%s", script.Name, script.Version, arch)
	l.Build = build
	l.BuildScript = script
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
	desk := &Desktop{Name: l.Name, Icon: l.Name + ".png", Executable: "AppRun", Categories: "Games"}
	if e := applyTempl(desktop, l.maniPath, desk); e != nil {
		return e
	}

	return nil
}

func (l *linux) postCompile() error {
	log.Println("packaging appimage...")
	if e := exec.Command("appimagetool-x86_64.AppImage", l.appDirPath).Run(); e != nil {
		return e
	}
	return nil
}

const desktop = `[Desktop Entry]
Name={{ .Name }}
Exec={{ .Executable }}
Icon={{ .Icon }}
Type=Application
Categories={{ .Categories }};`
