package animeworld

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"strconv"
)

type ArchiveMode int

type Archive struct {
	Mode       ArchiveMode
	Seasons    []*ArchiveSeason
	TotalPages int
	CurrentPage int
}

type ArchiveSeason struct {
	Url string
	Name string
	Plot string
	ImageUrl string
}

const (
	AzList ArchiveMode = iota
)

// NewArchive is a function that return
func NewArchive(mode ArchiveMode, page int) (*Archive, error) {

	a := Archive{
		Mode: mode,
		CurrentPage: page,
	}

	err := a.getInfo()
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// NextPage is a public method that gets all archive info from the next page
func (a *Archive) NextPage() error{

	if a.TotalPages != 0 && a.CurrentPage+1 > a.TotalPages {
		return errors.New("next page doesn't exist")
	}

	a.CurrentPage++
	err := a.getInfo()
	if err != nil {
		return err
	}

	return nil
}

// ForwardPage is a public method that gets all archive info from the forward page
func (a *Archive) ForwardPage() error{

	if a.CurrentPage-1 < 0 {
		return errors.New("forward page doesn't exist")
	}

	a.CurrentPage--
	err := a.getInfo()
	if err != nil {
		return err
	}

	return nil
}

// getInfo is a private method that gets all archive info
func (a *Archive) getInfo() error {
	c := colly.NewCollector()

	// get seasons info
	c.OnHTML("div.item", a.getSeasons)

	// get pages info
	c.OnHTML("div.paging-wrapper", a.getPagesInfo)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(a.makeArchiveUrl())
	if err != nil {
		return err
	}

	return nil
}

// getSeasons is a private method that gets base info for all series.
func (a *Archive) getSeasons(e *colly.HTMLElement) {

	as := ArchiveSeason{
		Url:      BaseUrl + e.ChildAttr("a", "href"),
		Name:     e.ChildText("a.name"),
		Plot:     e.ChildText("p"),
		ImageUrl: e.ChildAttr("img", "src"),
	}

	a.Seasons = append(a.Seasons, &as)
}

// getPagesInfo is a private method that gets current page and total pages.
func (a *Archive) getPagesInfo(e *colly.HTMLElement) {

	a.CurrentPage, _ = strconv.Atoi(e.ChildAttr("#page-input", "placeholder"))
	a.TotalPages, _ = strconv.Atoi(e.ChildText("span.total"))
}

// makeArchiveUrl is a private method that returns the correct url by archive mode.
func (a *Archive) makeArchiveUrl() string {

	// Check default value
	if a.CurrentPage == 0 {
		a.CurrentPage = 1
	}

	switch a.Mode {
	case AzList:
		return fmt.Sprintf("%s/az-list/?page=%d", BaseUrl, a.CurrentPage)
	default:
		return BaseUrl + "/az-list/"
	}
}