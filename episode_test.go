package animeworld

import "testing"

func TestEpisode_GetVideoUrl(t *testing.T) {
	
	e := Episode{
		Name:     "5",
		Url:      "https://www.animeworld.tv/play/3d-kanojo-real-girl.IdIRW/dWubp",
		VideoUrl: "",
	}

	err := e.GetVideoUrl()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("EPISODE:", e)
}


