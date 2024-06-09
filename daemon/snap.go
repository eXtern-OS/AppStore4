package daemon

import (
	"encoding/json"
	beatrix "github.com/eXtern-OS/Beatrix"
	"github.com/eXtern-OS/core9-common/models/app"
	"io"
	"net/http"
	"strings"
)

const snapBaseUrl = "https://api.snapcraft.io/v2/snaps/find?fields=media,description,publisher,title,version&q="

func runQuery(q string) ([]app.Snap, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", snapBaseUrl+q, nil)

	// Setting up required header
	req.Header.Set("Snap-Device-Series", "16")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var res app.SnapResults

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return res.Results, nil
}

func FetchSnap() ([]app.Snap, error) {
	var res []app.Snap
	for _, x := range strings.Split("abcdefghijklmnopqrstuvwxyz1234567890", "") {
		r, err := runQuery(x)
		if err != nil {
			go beatrix.SendError("Failed to fetch snap apps with query "+x, "AppStore4.Daemon.FetchSnap")
		}
		res = append(res, r...)
	}
	return res, nil
}
