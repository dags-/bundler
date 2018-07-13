package build

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func splitVersion(s string) []int {
	parts := strings.Split(s, ".")
	version := make([]int, 4)

	// parse the version points, pad right with zeros
	for i := 0; i < 4; i++ {
		if i < 3 && i < len(parts) {
			j, e := strconv.Atoi(parts[i])
			if e == nil {
				version[i] = j
				continue
			}
		}
		version[i] = 0
	}

	// try retrieve the build number from an environment variable
	if name, ok := os.LookupEnv("BUNDLE_BUILD_NUM"); ok {
		if number, ok := os.LookupEnv(name); ok {
			log.Println("build number:", number)
			build, e := strconv.Atoi(number)
			if e != nil {
				version[3] = build
			}
		} else {
			log.Println("no build number env variable set")
		}
	}

	return version
}
