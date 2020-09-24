package hltvscrape

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

const baseURL = "https://www.hltv.org"

func ExtractStats(statsPageURL string) (data MapData, err error) { // @TODO Work on the extract stats function.

	c := colly.NewCollector() // Could look to add options here for optimization
	return MapData{}, nil     //placeholder
}

func ExtractMatch(url string) (MatchData, error) {

	// Initilizate MatchData with the data we have from just the URL.
	match := MatchData{
		MatchURL: url,
		MatchID:  strings.Split(url, "/")[4], //Extracting match id from URl, must have https:// prefixed
	}

	// Now we start scraping into the Struct
	c := colly.NewCollector()
	// Extract Team Data for both teams via the 'team' div class(s)

	c.OnHTML(`.team`, func(e *colly.HTMLElement) {
		// Look for data of first team
		team0URL, exists := e.DOM.Find(".team1-gradient").Children().Attr("href") // Locates
		if exists {
			teamURLData(team0URL, &match.Team0)
			if whoWin(e) {
				match.Winner = 1
			}

		} //@TODO grab more data in both statements..
		team1URL, exists2 := e.DOM.Find(".team2-gradient").Children().Attr("href") // Locates
		if exists2 {
			teamURLData(team1URL, &match.Team1)
			if whoWin(e) {
				match.Winner = 2
			}

		}

		// Match.winner is set if we can find a 'won' div in their children.
		// If we cant find in either, match has to be a draw and defaults to that.

		// Team Data for both teams (from the match page)is now extracted.
		//selection.Find() // We can use this to search for elems/values
	})

	c.Visit(url)

	return match, nil
}

func whoWin(e *colly.HTMLElement) bool {
	s := e.DOM.Children().Find(".won")
	if len(s.Nodes) > 0 {
		return true
	}
	return false
}

func teamURLData(url string, t *Team) {
	t.TeamURL = baseURL + url
	t.TeamID = strings.Split(url, "/")[2]
	t.Name = strings.Split(url, "/")[3]

}

// @TODO abstract team data extraction into function. smilar to 'whoWin'
// @TODO check that match has starting via extracted unix timestamp of match start.
