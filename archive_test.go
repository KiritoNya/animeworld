package animeworld_test

import (
	"github.com/KiritoNya/animeworld"
	"testing"
)

func TestNewArchive(t *testing.T) {

	a, err := animeworld.NewArchive(animeworld.AzList, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ARCHIVE: ", a)
}

func TestArchive_NextPage(t *testing.T) {

	a, err := animeworld.NewArchive(animeworld.AzList, 1)
	if err != nil {
		t.Fatal(err)
	}

	err = a.NextPage()
	if err != nil {
		t.Fatal(err)
	}

	// test limit
	a.CurrentPage = a.TotalPages
	err = a.NextPage()
	if err == nil {
		t.Fatal("limit exceed")
	}

	t.Log("ARCHIVE:", a)
}

func TestArchive_ForwardPagePage(t *testing.T) {

	a, err := animeworld.NewArchive(animeworld.AzList, 2)
	if err != nil {
		t.Fatal(err)
	}

	err = a.ForwardPage()
	if err != nil {
		t.Fatal(err)
	}

	// test limit
	a.CurrentPage = 0
	err = a.ForwardPage()
	if err == nil {
		t.Fatal("limit exceed")
	}

	t.Log("ARCHIVE:", a)
}
