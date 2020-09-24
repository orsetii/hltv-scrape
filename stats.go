package hltvscrape

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

const baseURL = "https://www.hltv.org"

func ExtractStats(statsPageURL string) (data MapData, err error) { // @TODO Work on the extract stats function.

	//c := colly.NewCollector() // Could look to add options here for optimization
	return MapData{}, nil //placeholder
}

func ExtractMatch(url string) (match MatchData, err error) {

	// Initilizate MatchData with the data we have from just the URL.
	match = MatchData{
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
			res, score, err := whoWin(e)
			parseErr(err)
			match.Team0SeriesScore = score
			if res == 1 {
				match.Winner = 1
			}

		} //@TODO grab more data in both statements..
		team1URL, exists2 := e.DOM.Find(".team2-gradient").Children().Attr("href") // Locates
		if exists2 {
			teamURLData(team1URL, &match.Team1)
			res, score, err := whoWin(e)
			parseErr(err)
			match.Team1SeriesScore = score
			if res == 1 {
				match.Winner = 2
			}
		}
		if exists || exists2 {
			err = fmt.Errorf("couldn't get all teamdata")
		}

		// Match.winner is set if we can find a 'won' div in their children.
		// If we cant find in either, match has to be a draw and defaults to that.

		// Team Data for both teams (from the match page)is now extracted.
		//selection.Find() // We can use this to search for elems/values
	})
	if err != nil {
		return match, err
	}
	c.Visit(url)

	return match, nil
}

// whoWin extracts data from a 'teamx-gradient' html element from match pages.
// winner returns a 1 (for Won) that element has a 'won' div inside it. If there is a 'lost' div winner is a 0 (for lost), if a 'tie' div is found, winner is set to 2
func whoWin(e *colly.HTMLElement) (winner int8, score SeriesScore, err error) {
	s := e.DOM.Children().Find(".won")

	if len(s.Nodes) > 0 {
		// If we get here, there is a 'won' div in e's children.
		pscore, err := strconv.Atoi(s.Text())
		parseErr(err)
		score = SeriesScore(pscore)
		winner = 1
	} else if l := e.DOM.Children().Find(".lost"); len(l.Nodes) > 0 {
		// Getting here means no won div was found.
		// We now check for a 'lost' div, if that is not found, check for a
		// If we get here, there is a 'lost' div in e's children.
		winner = 0
		pscore, err := strconv.Atoi(l.Text())
		if err != nil {
			parseErr(err)
		}
		score = SeriesScore(pscore)
	} else if d := e.DOM.Children().Find(".tie"); len(d.Nodes) > 0 {
		winner = 2
		pscore, err := strconv.Atoi(s.Text())
		if err != nil {
			parseErr(err)
		}
		score = SeriesScore(pscore)
	} else {
		return 2, 0, fmt.Errorf("couldn't find a result in HTML")
	}
	return
}

func teamURLData(url string, t *Team) {
	t.TeamURL = baseURL + url
	t.TeamID = strings.Split(url, "/")[2]
	t.Name = strings.Split(url, "/")[3]

}

// @TODO abstract team data extraction into function. smilar to 'whoWin'
// @TODO check that match has starting via extracted unix timestamp of match start.

func parseErr(err error) error {

	if err != nil {
		return fmt.Errorf("error in extracting data from HTML: %s", err)
	}
	return nil

}
