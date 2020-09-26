package hltvscrape

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

const baseURL = "https://www.hltv.org"

// ExtractStats is used on a HLTV stats page, it extracts all data into applicable struct(s)
func ExtractStats(statsPageURL string) (data MapData, err error) { // @TODO Work on the extract stats function.

	//c := colly.NewCollector() // Could look to add options here for optimization
	return MapData{}, nil //placeholder
}

// ExtractMatch is used on a HLTV matchPage. It extracts all data from the matchpage, then will call functions to extract data from each player, and each map.
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
			extractTeamURLData(team0URL, &match.Team0)
			res, score, err := extractWinner(e)
			parseErr(err)
			match.Team0SeriesScore = score
			if res == 1 {
				match.Winner = 1
			}

		} //@TODO grab more data in both statements..
		team1URL, exists2 := e.DOM.Find(".team2-gradient").Children().Attr("href") // Locates
		if exists2 {
			extractTeamURLData(team1URL, &match.Team1)
			res, score, err := extractWinner(e)
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

	})
	// For each map link avaialable, get the statspage URL.
	c.OnHTML(`.results-stats`, func(e *colly.HTMLElement) {
		match.MapLinks = append(match.MapLinks, "https://www.hltv.org"+e.Attr("href"))
	})

	// This function extracts the length of the match (best of x)
	// It also extracts the stage of the tournament it is in.
	c.OnHTML(`.padding.preformatted-text`, func(e *colly.HTMLElement) {
		texts := strings.Split(e.Text, "\n")
		// We extract Best of X in top row.
		// Empty str in middle of slice
		// Match Context in bottom row
		if match.BestOfType, err = strconv.Atoi(strings.Split(texts[0], " ")[2]); err != nil {
			log.Printf("Error in extracting best of type. Attempted to extract from %s text.\n", e.Text)
			match.BestOfType = 0
		}
		match.Stage = texts[2][2:]
	})

	// This function extracts the name of the event that the match is played as a part of.
	c.OnHTML(`.event.text-ellipsis`, func(e *colly.HTMLElement) {
		match.Event = e.Text
		match.EventID = extractID(e.ChildAttr("a", "href"))
	})

	// This function extracts the exact unix time of the estimated match start
	c.OnHTML(`.timeAndEvent`, func(e *colly.HTMLElement) {
		match.MatchTimeEpoch, err = strconv.Atoi(e.ChildAttr(".time", "data-unix"))
		if err != nil {
			log.Printf("Could not get time from match page.\n")
		}
	})

	// Extract PickBans

	c.OnHTML(`.standard-box.veto-box`, func(e *colly.HTMLElement) {
		kids := e.DOM.Children()
		match.Vetos, err = extractVetos(kids.Children().Text())
		if err != nil {
			log.Println(err)
		}
	})

	if err != nil {
		return match, err
	}

	// We could extract total stats across maps however I do not want this data. Per map data is much more useful in data analysis.

	// Finally, we extract who played in the games. We could extract this from the teamdata, but obviously team rosters change, and individual performance can be extremely useful.
	// Counter to monitor which team we are looking at data for.
	var tCounter int = 0
	c.OnHTML(`.lineup.standard-box`, func(e *colly.HTMLElement) {
		var curTeam *Team
		tname := e.ChildAttr(".logo", "alt")
		log.Printf("Team: %s", tname)

		// Error Checking...
		if tname == "" {
			// If nothing in the ChildAttr must not be in right place, return out of function.
			return
		}
		if tCounter > 1 {
			// If we get here we have encountered an error, and are attempting to find data for a third team!
			log.Printf("Error parsing data, too many Team HTML elements \n")
			return
		}

		if tCounter == 0 {
			// We know we are working with team0:
			curTeam = &match.Team0

		} else if tCounter == 1 {
			curTeam = &match.Team1
		}

		// This iterates over each player dataframe for that team.
		e.ForEach(`.player.player-image`, func(playerIndex int, e *colly.HTMLElement) {
			log.Printf("Getting data for player %d of Team %s", playerIndex, tname)
			pData := Player{
				TeamPlayedFor: curTeam,
			}
			extractPlayerData(e.ChildAttr(`a`, "href"), &pData)
			curTeam.Players = append(curTeam.Players, pData)
		})
		tCounter++
		return
	})

	// All player data from match page extracted. Now, for each map played, extract map data
	for i, link := range match.MapLinks {
		match.MapsPlayed[i] = *ExtractMapData(link)
	}

	c.Visit(url)
	return match, nil
}

func extractPlayerData(url string, p *Player) {
	p.PlayerURL = baseURL + url
	p.PlayerID = extractID(url)
	p.Name = strings.Split(url, "/")[3]
}

func extractVetos(data string) (result VetoList, err error) {
	result = make(VetoList, 7)
	l := strings.Split(data, ".")
	if len(l) <= 1 {
		return VetoList{}, fmt.Errorf("could not extract Veto data")
	}
	log.Printf("%+#v", l)
	for i, v := range l {
		if i == 0 {
			continue
		}
		current := v[1 : len(v)-1]
		var pb int8
		splitVeto := strings.Split(current, " ")
		if strings.Contains(current, "removed") {
			// If we get here, we have hit a ban.
			pb = 1
		} else if strings.Contains(current, "picked") {
			// This func does nothing as picked value is 0, just here to check so not included in else statement.
		} else {
			// If we get here, there is either a fucked up string or this is the last map.
			// Check that this is the last element in the veto.
			if i != len(l)-1 {
				// If we get here, we have an error
				return result, fmt.Errorf("error extracting veto element %d", i)
			}

		}
		var mname string
		if i == len(l)-1 {
			pb = 2
			mname = splitVeto[0]
		} else {
			log.Println(current)
			mname = splitVeto[len(splitVeto)-1]
		}
		result[i-1] = veto{
			BanPick: pb,
			MapName: mname,
		}
	}
	return
}

// extractWinner extracts data from a 'teamx-gradient' html element from match pages.
// winner returns a 1 (for Won) that element has a 'won' div inside it. If there is a 'lost' div winner is a 0 (for lost), if a 'tie' div is found, winner is set to 2
func extractWinner(e *colly.HTMLElement) (winner int8, score SeriesScore, err error) {
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

func extractTeamURLData(url string, t *Team) {
	t.TeamURL = baseURL + url
	t.TeamID = extractID(url)
	t.Name = strings.Split(url, "/")[3]

}

// extracts ID from relative URL
func extractID(url string) (id string) {
	id = strings.Split(url, "/")[2]
	return
}

// ExtractMapData will extract all data needed for the MapData struct.
// @TODO complete this function!
func ExtractMapData(url string) (m *MapData) {
	return
}

func parseErr(err error) error {

	if err != nil {
		return fmt.Errorf("error in extracting data from HTML: %s", err)
	}
	return nil

}
