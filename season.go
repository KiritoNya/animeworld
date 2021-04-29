package animeworld

import (
	"errors"
	"fmt"
	"github.com/KiritoNya/htmlutils"
	"golang.org/x/net/html"
	"net/http"
	"time"
)

//Season is a struct with the information of season anime.
type Season struct {
	Cover         string
	Name          string
	Category      string
	Audio         string
	RelaseData    time.Time
	ReleaseSeason ReleasingSeason
	Vote          int
	Duration      int //Durata in minuti
	TotEpisode    int
	Visuals       int
	Trama         string
	Keywords      []string
	Relations     []*Season
	Suggestions   []*Season
	Episodes      []*Episode
	MyAnimeList   string
	Anilist       string
	Youtube       string
	node          *html.Node
	//Comments []Comment
	//Status State
	//Studio Studio
	//Genres []Genre
}

//ReleasingSeason is a struct with the information of a season. EX: "Inverno 2018"
type ReleasingSeason struct {
	Season string
	Year   int
}

//NewSeason is a constructor of Season object.
func NewSeason(urlSeason string) (*Season, error) {

	var s Season

	resp, err := http.Get(urlSeason)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	htmlBody, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	s.node = htmlBody

	return &s, nil
}

func (s *Season) GetName() error {

	info, err  := htmlutils.QuerySelector(s.node, "div", "class", "info col-md-9")
	if err != nil {
		return errors.New("Info of season not found")
	}

	nameHtml, err := htmlutils.QuerySelector(info[0], "h2", "class", "title")
	if err != nil {
		return errors.New("Title not found.")
	}
	s.Name = string(htmlutils.GetNodeText(nameHtml[0], "h2"))

	return nil
}

//GetAnilist extract anilist link.
func (s *Season) GetAnilist() error {

	node, err := htmlutils.QuerySelector(s.node, "a", "id", "anilist-button")
	if err != nil {
		return errors.New("Anilist button not found")
	}

	link, err := htmlutils.GetValueAttr(node[0], "a", "href")
	if err != nil {
		return errors.New("Error to extract Anilist link")
	}

	s.Anilist = string(link[0])

	return nil
}

//GetEpisodes create all episode object of season.
func (s *Season) GetEpisodes() error {

	nodes, err := htmlutils.QuerySelector(s.node, "li", "class", "episode")
	if err != nil {
		return errors.New("Episodes not found")
	}

	for numEp, node := range nodes {

		link, err := htmlutils.GetValueAttr(node, "a", "href")
		if err != nil {
			return errors.New(fmt.Sprintf("Link of episode %d not found", numEp))
		}

		ep, err := NewEpisode(BaseUrl + string(link[0]))
		if err != nil {
			return errors.New(fmt.Sprintf("Error to create new episode[%d] object.", numEp))
		}

		s.Episodes = append(s.Episodes, ep)

	}
	return nil
}
