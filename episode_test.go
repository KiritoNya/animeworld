package animeworld

import "testing"

func TestNewEpisode(t *testing.T) {
	_, err := NewEpisode("https://www.animeworld.tv/play/mushoku-tensei-isekai-ittara-honki-dasu.FwbNc/SKUFKi")
	if err != nil {
		t.Error("error to create object:", err)
	}
	t.Log("[OK]")
}

func TestNewRssAnimeworld2(t *testing.T) {

	result := "http://www.thehiddendungeononlyicanenter.net/download-file.php?id=DLL/ANIME/MushokuTensei/MushokuTensei_Ep_01-02_SUB_ITA_v0.mp4"

	ep, err := NewEpisode("https://www.animeworld.tv/play/mushoku-tensei-isekai-ittara-honki-dasu.FwbNc/SKUFKi")
	if err != nil {
		t.Error("error to create object:", err)
	}

	err = ep.GetDirectLink()
	if err != nil {
		t.Error("error to obtain direct link:", err)
	}

	if ep.DirectLink != result {
		t.Error("not obtain", result, "but obtain", ep.DirectLink)
	} else {
		t.Log("[OK]")
	}
}