package animeworld

import (
	"testing"
)

func TestNewArchive(t *testing.T) {
	_, err := NewArchive("all")
	if err != nil {
		t.Fatal("Error to create archive object")
	}
	t.Log("[OK] Create object")
}

func TestArchive_GetSeason(t *testing.T) {

	a, err := NewArchive("e")
	if err != nil {
		t.Fatal("Error to create archive object")
	}

	err = a.GetSeason()
	if err != nil {
		t.Fatal(err)
	}

	for _, season := range a.Seasons {

		err := season.GetName()
		if err != nil {
			t.Fatal(err)
		}

		t.Log(season.Name)

	}
}
