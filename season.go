package animeworld

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

//Season is a struct with the information of season anime.
type Season struct {
	Cover         string
	Title          string
	Category      string
	Audio         string
	ReleaseData    string
	ReleaseSeason string
	Vote          float64
	Duration      int //Durata in minuti
	TotEpisode    int
	Visuals       int
	Plot         string
	Keywords      []string
	Relations     []*SeasonRelated
	Suggestions   []*SeasonSuggestion
	Episodes      []*Episode
	MyAnimeList   string
	Anilist       string
	Mangaworld	  string
	Youtube       string
	node          *html.Node
	//Comments []Comment
	Status string
	Studio string
	Genres []string
}

// SeasonSuggestion is a struct with the information of suggestion season.
type SeasonSuggestion struct {
	Title string
	Url string
	Cover string
}

// SeasonRelated is a struct with the information of related season.
type SeasonRelated struct {
	Title string
	Url string
	Cover string
	Category string
	Year int
	Duration int
}


// NewSeason is a Season constructor that return the season with all its info
func NewSeason(url string) (*Season, error) {

	var s Season

	c := colly.NewCollector()

	// get seasons info
	c.OnHTML("#thumbnail-watch", s.getCover)

	// get pages info
	c.OnHTML("div.widget.info", s.getInfo)

	// get keywords
	c.OnHTML("#tagsReload", s.getKeywords)

	// get related
	c.OnHTML("div.widget-body.related", s.getRelated)

	// get suggestions
	c.OnHTML("div.film-list.interesting", s.getSuggestions)

	// get links
	c.OnHTML("#controls", s.getSeasonLinks)

	// get episodes
	c.OnHTML("div.server.active", s.getEpisodes)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// getCover is a private method that gets the season cover
func (s *Season) getCover(e *colly.HTMLElement) {
	s.Cover = e.ChildAttr("img", "src")
}

// getInfo is a private method that gets the basic info
func (s *Season) getInfo(e *colly.HTMLElement) {

	s.Title = e.ChildText("h2.title")
	e.ForEach("dd", func (position int, dd *colly.HTMLElement) {
		switch position {
		case 0:
			s.Category = strings.ReplaceAll(dd.Text, "  ", "")
			s.Category = strings.ReplaceAll(s.Category, "\n", "")
		case 1:
			s.Audio = dd.Text
		case 2:
			s.ReleaseData = dd.Text
		case 3:
			s.ReleaseSeason = dd.Text
		case 4:
			s.Studio = strings.ReplaceAll(dd.Text, "\n", "")
			s.Studio = strings.ReplaceAll(s.Studio, "  ", "")
		case 5:
			genres := strings.ReplaceAll(dd.Text, "\n", "")
			genres = strings.ReplaceAll(genres, "  ", "")
			s.Genres = strings.Split(genres, ",")
		case 6:
			s.Vote, _ = strconv.ParseFloat(dd.ChildText("#average-vote"), 32)
		case 7:
			s.Duration, _ = strconv.Atoi(strings.Split(dd.Text, " ")[0])
		case 8:
			s.TotEpisode, _ = strconv.Atoi(dd.Text)
		case 9:
			s.Status = strings.ReplaceAll(dd.Text, "\n", "")
			s.Status = strings.ReplaceAll(s.Status, "  ", "")
		case 10:
			s.Visuals, _ = strconv.Atoi(strings.ReplaceAll(dd.Text, ".", ""))
		}
	})
	s.Plot = e.ChildText("div.desc")
}

// getKeywords is a private method that gets season keywords
func (s *Season) getKeywords(e *colly.HTMLElement) {
	keywords := strings.ReplaceAll(e.Text, "  ", "")
	s.Keywords = strings.Split(keywords, " - ")
}

// getSuggestions is a private method that gets suggestion season
func (s *Season) getSuggestions(e *colly.HTMLElement) {

	e.ForEach("div.item", func (_ int, ele *colly.HTMLElement){
		var sr SeasonSuggestion
		sr.Title = e.ChildAttr("a.name", "title")
		sr.Cover = e.ChildAttr("img", "src")
		sr.Url = BaseUrl + e.ChildAttr("a.name", "href")
		s.Suggestions = append(s.Suggestions, &sr)
	})
}

// getRelated is a private method that gets related season
func (s *Season) getRelated(e *colly.HTMLElement) {

	e.ForEach("div.item", func (_ int, ele *colly.HTMLElement){
		var sr SeasonRelated
		sr.Title = e.ChildText("a.name")
		sr.Cover = e.ChildAttr("img", "src")
		sr.Url = BaseUrl + e.ChildAttr("a.name", "href")

		// gets season info
		info := strings.Split(e.ChildText("p"), " - ")
		sr.Category = info[0]
		sr.Year, _ = strconv.Atoi(info[1])
		sr.Duration, _ = strconv.Atoi(strings.Split(info[2], " ")[0])

		s.Relations = append(s.Relations, &sr)
	})
}

// getSeasonLinks is a private method that gets season link like MAL, anilist or youtube trailer
func (s *Season) getSeasonLinks(e *colly.HTMLElement) {

	s.MyAnimeList = e.ChildAttr("#mal-button", "href")
	s.Anilist = e.ChildAttr("#anilist-button", "href")
	s.Mangaworld = e.ChildAttr("#mangaworld-button", "href")
	s.Youtube = e.ChildAttr("div.trailer.control.tip.tippy-desktop-only", "data-url")
}

// getEpisodes is a private method that gets season episodes
func (s *Season) getEpisodes(e *colly.HTMLElement) {

	e.ForEach("li.episode", func (_ int, ele *colly.HTMLElement) {
		var ep Episode
		ep.Url = BaseUrl + ele.ChildAttr("a", "href")
		ep.Name = ele.ChildText("a")

		s.Episodes = append(s.Episodes, &ep)
	})
}