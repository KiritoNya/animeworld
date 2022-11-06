package animeworld

import (
	"encoding/json"
	"testing"
)

func TestNewSeason(t *testing.T) {

	s, err := NewSeason("https://www.animeworld.tv/play/3d-kanojo-real-girl.IdIRW/4g_EM")
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.MarshalIndent(s, " ", "\t")

	t.Log("SEASON:", string(data))
}
