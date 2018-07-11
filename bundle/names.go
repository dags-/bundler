package bundle

func toNormal(name string) string {
	switch name {
	case "darwin":
		return "macOS"
	case "amd64":
		return "x64"
	case "386":
		return "x32"
	}
	return name
}
