package animeworld

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KiritoNya/htmlutils"
	"golang.org/x/net/html"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type Episode struct {
	Url       string
	StreamLinks map[string]string
	Number     []float64
	DownloadLink string //TODO: Gestirlo in futuro per rendere il pacchetto generale
	htmlPage *html.Node
}

type Server struct {
	ID int
	Name string
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

	return &Episode{Url: link, htmlPage: htmlBody, StreamLinks: make(map[string]string)}, nil
}

//GetNumber extract the number of episode.
func (ep *Episode) GetNumber() error {

	nodes, err := htmlutils.QuerySelector(ep.htmlPage, "li", "class", "episode")
	if err != nil {
		return errors.New("Episode not found.")
	}

	//Foreach episode
	for _, node := range nodes {

		link, err := htmlutils.QuerySelector(node, "a", "class", "active")
		if err == nil {

			num, err := htmlutils.GetValueAttr(link[0], "a", "data-episode-num")
			if err != nil {
				return errors.New("Episode active number not found.")
			}
			numString := string(num[0])

			episodes := strings.Split(numString, "-")

			for _, episode := range episodes {

				numFloat, err := strconv.ParseFloat(episode, 64)
				if err != nil {
					return err
				}

				ep.Number = append(ep.Number, numFloat)
			}
		}
	}

	if err != nil {
		return errors.New("Error to get active episode.")
	}

	return nil
}

//GetStreamLinks extract links of episodes di vari server.
func (ep *Episode) GetStreamLinks() error {

	if ep.Number == nil {
		err := ep.GetNumber()
		if err != nil {
			return errors.New("Error to ge number: " + err.Error())
		}
	}

	servers, err := getServers(ep.htmlPage)
	if err != nil {
		return err
	}

	//Get section of servers episodes
	serversHtml, err := htmlutils.QuerySelector(ep.htmlPage, "div", "class", "server")
	if err != nil {
		return errors.New("Error to get servers sections.")
	}

	for _, server := range servers {

		var episodesHtml []*html.Node

		//Get server episodes html list
		for _, serverHtml := range serversHtml {

			node, err := htmlutils.QuerySelector(serverHtml, "div", "data-id", strconv.Itoa(server.ID))
			if err != nil {
				return errors.New("Server with id " + strconv.Itoa(server.ID) + " not found")
			}

			//Get all episodes
			episodesHtml, err = htmlutils.QuerySelector(node[0], "li", "class", "episode")
			if err != nil {
				return errors.New("Episode not found.")
			}

		}

		for i, episodeHtml := range episodesHtml {

			var episodeNums []float64

			link, err := htmlutils.GetGeneralTags(episodeHtml, "a")
			if err != nil {
				return errors.New(fmt.Sprint("Tag a of episode ", i, " not found!"))
			}

			num, err := htmlutils.GetValueAttr(link[0], "a", "data-episode-num")
			if err != nil {
				return errors.New("Episode active number not found.")
			}
			numString := string(num[0])

			episodes := strings.Split(numString, "-")

			for _, episode := range episodes {

				numFloat, err := strconv.ParseFloat(episode, 64)
				if err != nil {
					return err
				}

				episodeNums = append(episodeNums, numFloat)
			}

			//Check if is the same episodes
			if equalSliceFloat(ep.Number, episodeNums) {

				var objmap map[string]json.RawMessage
				var directLink string

				//Extract direct link
				href, err := htmlutils.GetValueAttr(link[0], "a", "href")
				if err != nil {
					return errors.New("Link of episode not found")
				}

				resp, err := http.Get(episodeInfoApi+filepath.Base(string(href[0])))
				if err != nil {
					return errors.New("Error to get response of episode info api.")
				}
				defer resp.Body.Close()

				err = json.NewDecoder(resp.Body).Decode(&objmap)
				if err != nil {
					return errors.New("Error to decode response of episode info api.")
				}

				err = json.Unmarshal(objmap["grabber"], &directLink)
				if err != nil {
					return errors.New("Error to decode response of episode info api.")
				}

				ep.StreamLinks[server.Name] = directLink

				fmt.Println(directLink)
			}
		}
	}
	return nil
}

//GetDirectLink add direct link at the Episode object. EX: .mp4
func (ep *Episode) GetDownloadLink() error {
	tagsA, err := htmlutils.QuerySelector(ep.htmlPage, "a", "id", "downloadLink")
	if err != nil {
		return err
	}

	links, err := htmlutils.GetValueAttr(tagsA[0], "a", "href")
	if err != nil {
		return err
	}

	if string(links[0]) != "" {
		ep.DownloadLink = strings.Replace(string(links[0]), "download-file.php?id=", "", -1)
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

		ep.DownloadLink = strings.Replace(string(links[0]), "download-file.php?id=", "", -1)
	}
	return  nil
}

