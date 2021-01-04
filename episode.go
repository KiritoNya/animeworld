package animeworld

import (
	"errors"
	"github.com/KiritoNya/htmlutils"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type Episode struct {
	Link       string
	Number     []float64
	DirectLink string
	htmlPage *html.Node
}

//NewEpisode is an constructor of Episode object.
func NewEpisode(link string) (*Episode, error) {

	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	htmlBody, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Episode{Link: link, htmlPage: htmlBody }, nil
}

//GetDirectLink add direct link at the Episode object. EX: .mp4
func (ep *Episode) GetDirectLink() error {
	tagsA, err := htmlutils.QuerySelector(ep.htmlPage, "a", "id", "downloadLink")
	if err != nil {
		return err
	}

	links, err := htmlutils.GetValueAttr(tagsA[0], "a", "href")
	if err != nil {
		return err
	}

	if string(links[0]) != "" {
		ep.DirectLink = strings.Replace(string(links[0]), "download-file.php?id=", "", -1)
	} else {
		tagsA, err = htmlutils.QuerySelector(ep.htmlPage, "a", "id", "alternativeDownloadLink")
		if err != nil {
			return err
		}

		links, err = htmlutils.GetValueAttr(tagsA[0], "a", "href")
		if err != nil {
			return err
		}

		if string(links[0]) == "" {
			return errors.New("Direct alternative link not found!")
		}

		ep.DirectLink = strings.Replace(string(links[0]), "download-file.php?id=", "", -1)
	}
	return  nil
}