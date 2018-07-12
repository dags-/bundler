package build

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type windows struct {
	*Build
	*Script
	exePath string
}

func (w *windows) artifact() string {
	return w.exePath
}

func (w *windows) executable() string {
	return w.exePath
}

func (w *windows) init(script *Script, build *Build, arch string) {
	name := fmt.Sprintf("%s-%s-%s.exe", script.Name, script.Version, arch)
	w.Build = build
	w.Script = script
	w.exePath = filepath.Join(script.Output, "windows", name)
}

func (w *windows) preCompile() error {
	log.Println("writing version-info...")
	f, e := os.Create("versioninfo.json")
	if e != nil {
		return e
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(w.manifest())
}

func (w *windows) postCompile() error {
	log.Println("cleaning up...")
	os.Remove("resource.syso")
	os.Remove("versioninfo.json")
	return nil
}

func (w *windows) manifest() interface{} {
	version := splitVersion(w.Version)
	return &VersionInfo{
		FixedFileInfo: FixedVersionInfo{
			FileVersion: Version{
				Major: version[0],
				Minor: version[1],
				Patch: version[2],
				Build: version[3],
			},
			ProductVersion: Version{
				Major: version[0],
				Minor: version[1],
				Patch: version[2],
				Build: version[3],
			},
			FileFlagsMask: "3f",
			FileFlags:     "00",
			FileOS:        "040004",
			FileType:      "01",
			FileSubType:   "00",
		},
		StringFileInfo: StringFileInfo{
			ProductName:    w.Name,
			ProductVersion: w.Version,
		},
		VarFileInfo: VarFileInfo{
			Translation: Translation{
				LangID:    "0409",
				CharsetID: "04B0",
			},
		},
		IconPath:     w.Build.Icon,
		ManifestPath: "",
	}
}

type VersionInfo struct {
	FixedFileInfo  FixedVersionInfo
	StringFileInfo StringFileInfo
	VarFileInfo    VarFileInfo
	IconPath       string
	ManifestPath   string
}

type FixedVersionInfo struct {
	FileVersion    Version
	ProductVersion Version
	FileFlagsMask  string
	FileFlags      string
	FileOS         string
	FileType       string
	FileSubType    string
}

type StringFileInfo struct {
	Comments         string
	CompanyName      string
	FileDescription  string
	FileVersion      string
	InternalName     string
	LegalCopyright   string
	LegalTrademarks  string
	OriginalFilename string
	PrivateBuild     string
	ProductName      string
	ProductVersion   string
	SpecialBuild     string
}

type VarFileInfo struct {
	Translation Translation
}

type Translation struct {
	LangID    string
	CharsetID string
}

type Version struct {
	Major int
	Minor int
	Patch int
	Build int
}
