package build

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type linux struct {
	*Build
	*Script
	arch       string
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
	l.arch = arch
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
	if runtime.GOOS != "linux" {
		return nil
	}

	log.Println("packaging appimage...")
	tool, e := l.getImageTool()
	if e != nil {
		return e
	}

	dir, imgTool := filepath.Split(tool)
	_, appDir := filepath.Split(l.appDirPath)

	c := exec.Command(imgTool, appDir)
	c.Dir = dir

	if e := c.Run(); e != nil {
		return e
	}

	return nil
}

func (l *linux) getImageTool() (string, error) {
	name := fmt.Sprintf("appimagetool-%s.AppImage", l.arch)
	path := filepath.Join(l.Output, "linux", name)
	if exists(path) {
		return path, nil
	}

	r, e := getLatestRelease("AppImage", "AppImageKit")
	if e != nil {
		return "", e
	}

	asset := "appimagetool-i686.AppImage"
	if toNormal(l.arch) == "x64" {
		asset = "appimagetool-x86_64.AppImage"
	}

	a, e := r.findAsset(asset)
	if e != nil {
		return "", e
	}

	f, e := download(a.Download, path)
	if e != nil {
		return "", e
	}
	defer f.Close()

	return path, f.Chmod(os.ModePerm)
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
