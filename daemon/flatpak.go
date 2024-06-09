package daemon

import (
	"encoding/json"
	"github.com/eXtern-OS/core9-common/models/app"
	"io"
	"net/http"
)

const flatpakBaseURL = "https://flathub.org/api/v1/apps"

func FetchFlatpak() ([]app.Flatpak, error) {
	resp, err := http.Get(flatpakBaseURL)

	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res []app.Flatpak
	return res, json.Unmarshal(b, &res)
}
