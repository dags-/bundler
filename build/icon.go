package build

import (
	"image"
	"io"
	"os"

	"github.com/Kodeworks/golang-image-ico"
	"github.com/tonyhb/goicns"
)

type converter func(io.WriteCloser, image.Image) error

func icon(script *Script, platform string) string {
	switch platform {
	case "darwin":
		return script.macIcon
	case "windows":
		return script.winIcon
	default:
		return script.Icon
	}
}

func loadImage(path string) (image.Image, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	img, _, e := image.Decode(f)
	return img, e
}

func writeIcon(c converter, img image.Image, path string) (string, error) {
	mustFile(path)
	out, e := os.Create(path)
	if e != nil {
		return "", nil
	}
	defer out.Close()
	return path, c(out, img)
}

func convertIco(w io.WriteCloser, img image.Image) error {
	return ico.Encode(w, img)
}

func convertIcns(w io.WriteCloser, img image.Image) error {
	icn := goicns.NewICNS(img)
	e := icn.Construct()
	if e != nil {
		return e
	}
	_, e = icn.WriteTo(w)
	return e
}
