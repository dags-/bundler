package bundle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type windows struct{}

func (w *windows) Artifact(b *BuildScript, p *Build, arch string) string {
	name := fmt.Sprintf("%s-%s-%s.exe", b.Name, b.Version, toNormal(arch))
	return filepath.Join(b.Output, "windows", name)
}

func (w *windows) ExecPath(b *BuildScript, p *Build, arch string) string {
	return w.Artifact(b, p, arch)
}

func (w *windows) WriteIcon(b *BuildScript, p *Build, arch string) error {
	return nil
}

func (w *windows) WriteManifest(b *BuildScript, p *Build, arch string) error {
	version := ver(b.Version)

	manifest := &VersionInfo{
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
			ProductName:    b.Name,
			ProductVersion: b.Version,
		},
		VarFileInfo: VarFileInfo{
			Translation: Translation{
				LangID:    "0409",
				CharsetID: "04B0",
			},
		},
		IconPath:     p.Icon,
		ManifestPath: "",
	}

	f, e := os.Open("versioninfo.json")
	if e != nil {
		f, e = os.Create("versioninfo.json")
		if e != nil {
			return e
		}
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(manifest)
}

func ver(s string) []int {
	parts := strings.Split(s, ".")
	ver := make([]int, 4)
	for i := 0; i < 4; i++ {
		val := 0
		if i < len(parts) {
			j, e := strconv.Atoi(parts[i])
			if e == nil {
				val = j
			}
		}
		ver[i] = val
	}
	return ver
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
