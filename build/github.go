package build

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Release struct {
	Assets []*Asset `json:"assets"`
}

type Asset struct {
	Name     string `json:"name"`
	Download string `json:"browser_download_url"`
}

func getLatestRelease(user, repo string) (*Release, error) {
	url := fmt.Sprint("https://api.github.com/repos/", user, "/", repo, "/releases/latest")
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var release Release
	e = json.NewDecoder(r.Body).Decode(&release)
	if e != nil {
		return nil, e
	}
	return &release, nil
}

func (r *Release) findAsset(name string) (*Asset, error) {
	for _, a := range r.Assets {
		if a.Name == name {
			return a, nil
		}
	}
	return nil, errors.New("no match")
}
