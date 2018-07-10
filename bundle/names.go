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

func toInternal(name string) string {
	switch name {
	case "osx":
		fallthrough
	case "macOS":
		return "darwin"
	case "win":
		return "windows"
	case "x64":
		return "amd64"
	case "x32":
		return "386"
	}
	return name
}
