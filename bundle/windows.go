package bundle

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/josephspurrier/goversioninfo"
)

type windows struct{}

func (w *windows) ExecPath(b *Build, p *Platform, arch string) string {
	name := fmt.Sprintf("%s-%s-%s.exe", b.Name, b.Version, arch)
	return filepath.Join(b.Output, "windows", name)
}

func (w *windows) WriteIcon(b *Build, p *Platform, arch string) error {
	return nil
}

func (w *windows) WriteManifest(b *Build, p *Platform, arch string) error {
	version := ver(b.Version)

	manifest := &goversioninfo.VersionInfo{
		FixedFileInfo: goversioninfo.FixedFileInfo{
			FileVersion: goversioninfo.FileVersion{
				Major: version[0],
				Minor: version[1],
				Patch: version[2],
				Build: version[3],
			},
			ProductVersion: goversioninfo.FileVersion{
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
		StringFileInfo: goversioninfo.StringFileInfo{
			ProductName:    b.Name,
			ProductVersion: b.Version,
		},
		VarFileInfo: goversioninfo.VarFileInfo{
			Translation: goversioninfo.Translation{
				LangID:    goversioninfo.LngUKEnglish,
				CharsetID: goversioninfo.CsUnicode,
			},
		},
		IconPath:     p.Icon,
		ManifestPath: "",
	}

	f, e := os.Open("versioninfo.json")
	if e != nil {
		log.Println(e)
		f, e = os.Create("versioninfo.json")
		if e != nil {
			log.Println(e)
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
