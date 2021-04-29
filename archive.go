package animeworld

import (
	"errors"
	"github.com/KiritoNya/htmlutils"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
)

type Archive struct {
	Seasons []*Season
	TotalPages int
	url string
	htmlPage *html.Node
}

func NewArchive(method string) (*Archive, error) {

	var a Archive
	var err error

	a.url = composeURL(method)
	a.htmlPage, err = doRequest(a.url)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Archive) GetSeason() (err error) {

	if a.htmlPage == nil {
		a.htmlPage, err = doRequest(composeURL("all"))
		if err != nil {
			return err
		}
	}

	err = a.GetTotalPages()
	if err != nil {
		return err
	}

	for i:=0; i<a.TotalPages ; i++ {

		link := a.url + "?page=" + strconv.Itoa(i)

		htmlPage, err := doRequest(link)
		if err != nil {
			return err
		}

		//Get container
		container, err := htmlutils.QuerySelector(htmlPage, "div", "class", "widget az-list")
		if err != nil {
			return errors.New("Container not found!")
		}

		//Get all seasons html sections
		seasonsNode, err := htmlutils.QuerySelector(container[0], "div", "class", "item")
		if err != nil {
			return errors.New("Error to obtain seasons section in the archive")
		}

		for i, seasonNode := range seasonsNode {

			if i == 0 {
				continue
			}

			tagA, err := htmlutils.GetGeneralTags(seasonNode, "a")
			if err != nil {
				return errors.New("Error to obtain tag \"a\" of seasons section")
			}

			href, err := htmlutils.GetValueAttr(tagA[0], "a", "href")
			if err != nil {
				return errors.New("Error to obtain link of seasons section")
			}

			season, err := NewSeason( BaseUrl + string(href[0]))
			if err != nil {
				return errors.New("Error to create season object")
			}

			a.Seasons = append(a.Seasons, season)
		}

	}

	return nil
}

func (a *Archive) GetTotalPages() error {
	div, err := htmlutils.QuerySelector(a.htmlPage, "div", "class", "paging-wrapper")
	if err != nil {
		return errors.New("Error to get paging wrapper for calculate max pages")
	}

	totalNum, err := htmlutils.QuerySelector(div[0], "span", "class", "total")
	if err != nil {
		return errors.New("Error to get total number of page.")
	}

	numTotal := htmlutils.GetNodeText(totalNum[0], "span")

	a.TotalPages, err = strconv.Atoi(string(numTotal))
	if err != nil {
		return err
	}

	return nil
}

func doRequest(link string) (*html.Node, error) {

	resp, err := http.Get(link)
	if err != nil {
		return nil, errors.New("Errors to get archive page: " + err.Error())
	}
	defer resp.Body.Close()

	htmlPage, err := html.Parse(resp.Body)
	if err != nil {
		return nil, errors.New("Error to parse archive page: " + err.Error())
	}

	return htmlPage, nil
}

func composeURL(method string) string {
	var link string

	if method == "all" {
		link = ArchiveUrl
	} else {
		link= ArchiveUrl + method
	}

	return link
}

