package build

import (
	"fmt"
	"io/ioutil"
	"os"
)

var contents = `package metainfo

const NAME = "%s"
const VERSION = "%s"
`

func generateInfo(name, version string) error {
	e := os.Mkdir("metainfo", os.ModePerm)
	if !os.IsExist(e) {
		return e
	}
	body := fmt.Sprintf(contents, name, version)
	return ioutil.WriteFile("metainfo/metainfo.go", []byte(body), os.ModePerm)
}
