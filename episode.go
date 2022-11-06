package animeworld

import (
	"encoding/json"
	"net/http"
	"path/filepath"
)

// Episode is a struct with the information of episode.
type Episode struct {
	Name string
	Url string
	VideoUrl string
}

func (ep *Episode) GetVideoUrl() error {

	// get episode id
	epID := filepath.Base(ep.Url)

	// do request
	resp, err := http.Get(episodeInfoApi + epID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse json
	var raw map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return err
	}

	ep.VideoUrl = raw["grabber"].(string)
	return nil
}