package hltvscrape

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const baseURL = "https://www.hltv.org"

// ExtractStats is used on a HLTV stats page, it extracts all data into applicable struct(s)
func ExtractStats(url string) (data *MapData) {
	data = new(MapData)
	data.statPageURL = url
	log.Println("Loading stats from: ", url)
	c := colly.NewCollector() // Could look to add options here for optimization
	var tCounter int
	c.OnHTML(`.stats-table`, func(e *colly.HTMLElement) {
		// Delve into each table for each team.
		// Check what team we are in
		if tCounter != 0 && tCounter != 1 {
			return
		}

		// Get Player data in table..
		e.ForEach(`tbody`, func(i int, e *colly.HTMLElement) {
			// For each row in table...
			e.ForEach(`tr`, func(i int, e *colly.HTMLElement) {
				var err error
				var p *PlayerMapPerf
				if tCounter == 0 {
					p = &data.Team0PlayerData[i]
				} else if tCounter == 1 {
					p = &data.Team1PlayerData[i]
				} else {
					log.Printf("Team index overrun at L37 in stats.go. Attempting to load teamid %d", tCounter)
				}
				// Get Player Data
				var ok bool
				plink, ok := e.DOM.Find(`.st-player`).Find("a").Attr("href")
				if !ok {
					log.Printf("Error in getting data for player id %d of Team%d\n", i, tCounter)
				}
				// Get name from URL
				splitLink := strings.Split(plink, "/")
				p.Name = splitLink[len(splitLink)-1]
				// Extract kills and headshots
				p.Kills, p.Headshots, err = splitValSubVal(e.DOM.Find(`.st-kills`).Text())
				if err != nil {
					log.Println("Could not get kill/headshot count. Error: ", err)
				}
				p.Assists, p.FlashAssists, err = splitValSubVal(e.DOM.Find(`.st-assists`).Text())
				if err != nil {
					log.Println("Could not get assist/flash assist count. Error: ", err)
				}

				p.Deaths, err = strconv.Atoi(e.DOM.Find(`.st-deaths`).Text())
				kastp := e.DOM.Find(`.st-kdratio`).Text()[:4]
				p.KASTPercentage, err = parseFloat(kastp)
				if err != nil {
					log.Println("Could not get KAST percentage. Error: ", err)
				}
				var pkddif string
				pkddif = e.DOM.Find(`.st-kddiff`).Text()
				if pkddif == "" {
					log.Printf("Error in getting first kill difference. No data found.")
				} else {
					p.KillDeathDiff, err = parseDiff(pkddif)
				}
				if err != nil {
					log.Println("Error getting KD Diff into int. Error: ", err)
				}
				p.ADR, err = parseFloat(e.DOM.Find(`.st-adr`).Text())
				if err != nil {
					log.Println("Error getting ADR into float value. Error: ", err)
				}
				// fk diff is players first kills - first kills on the player
				var pfkdiff string
				pfkdiff = e.DOM.Find(`.st-fkdiff`).Text()
				if pfkdiff != "" {
					p.FirstKillsDiff, err = parseDiff(pfkdiff)
				} else {
					log.Printf("Error in getting first kill difference. No data found.")
				}

				p.Rating, err = parseFloat(e.DOM.Find(`.st-rating`).Text())
				if err != nil {
					log.Printf("Error getting and parsing rating for player %d. Error: %s\n", i, err)
				}
				log.Printf("Player Data extracted: %+v", *p)
			})

		})
		tCounter++
	})
	c.OnHTML(`.match-info-row`, func(e *colly.HTMLElement) {
		var err error
		childs := e.DOM.Children().Text()
		// We extract Text from the each row.
		if strings.Contains(childs, "Breakdown") {
			// If we are looking at the breakdown row
			splitbr := Splitter(childs, " :()")
			// [0] = Team0 Total
			// [1] = Team1 Total
			// [2] = Team0 First Half
			// [3] = Team1 First Half
			// [4] = Team0 Second Half
			// [5] = Team1 Second Half

			// Extract Team0 Total Score
			data.Team0ScoreTotal, err = strconv.Atoi(splitbr[0])
			if err != nil {
				log.Printf("Could not get Team0 Total Score. Error: %s\n", err)
			}
			// Extract Team1 Total Score
			data.Team1ScoreTotal, err = strconv.Atoi(splitbr[1])
			if err != nil {
				log.Printf("Could not get Team1 Total Score. Error: %s\n", err)
			}
			// Extract Team0 First Half Score
			data.Team0ScoreFirstHalf, err = strconv.Atoi(splitbr[2])
			if err != nil {
				log.Printf("Could not get Team0 First Half Score. Error: %s\n", err)
			}
			// Extract Team1 First Half Score
			data.Team1ScoreFirstHalf, err = strconv.Atoi(splitbr[3])
			if err != nil {
				log.Printf("Could not get Team0 Sec Half Score. Error: %s\n", err)
			}
			// Extract Team0 Second Half Score
			data.Team0ScoreSecondHalf, err = strconv.Atoi(splitbr[4])
			if err != nil {
				log.Printf("Could not get Team0 First Half Score")
			}
			// Extract Team1 Second Half Score
			data.Team1ScoreSecondHalf, err = strconv.Atoi(splitbr[5])
			if err != nil {
				log.Printf("Could not get Team0 First Half Score")
			}

		} else if strings.Contains(childs, "Team rating") {
			pts := Splitter(childs, " :")
			data.Team0TeamRating, err = parseFloat(pts[0][:4])
			if err != nil {
				log.Printf("Could not get Team0 Team Rating. Error: %s\n", err)
			}
			data.Team1TeamRating, err = parseFloat(pts[1][:4])
			if err != nil {
				log.Printf("Could not get Team1 Team Rating. Error: %s\n", err)
			}
		} else if strings.Contains(childs, "First kills") {
			pts := Splitter(childs, " :F")
			data.Team0FirstKills, err = strconv.Atoi(pts[0])
			if err != nil {
				log.Printf("Could not get Team0 First Kills. Error: %s\n", err)
			}
			data.Team1FirstKills, err = strconv.Atoi(pts[1])
			if err != nil {
				log.Printf("Could not get Team1 First Kills. Error: %s\n", err)
			}
		} else if strings.Contains(childs, "Clutches won") {
			pts := Splitter(childs, " :C")
			data.Team0FirstKills, err = strconv.Atoi(pts[0])
			if err != nil {
				log.Printf("Could not get Team0 Clutches Won. Error: %s\n", err)
			}
			data.Team1ClutchesWon, err = strconv.Atoi(pts[1])
			if err != nil {
				log.Printf("Could not get Team1 Clutches Won. Error: %s\n", err)
			}
		}
	})
	if data.Team0ScoreTotal > data.Team1ScoreTotal {
		data.Winner = 1
	} else if data.Team0ScoreTotal < data.Team1ScoreTotal {
		data.Winner = 2
	}
	c.OnHTML(`.match-info-box`, func(e *colly.HTMLElement) {
		for _, v := range Maps {
			if strings.Contains(e.Text, v) {
				data.MapName = v
				break
			}
		}
		if data.MapName == "" {
			log.Printf("Could not find a map in the HTML.")
		}
	})
	c.Visit(url)

	return
}

// Splitter splits s string by all runes in the splits string
func Splitter(s string, splits string) []string {
	m := make(map[rune]int)
	for _, r := range splits {
		m[r] = 1
	}

	splitter := func(r rune) bool {
		return m[r] == 1
	}

	return strings.FieldsFunc(s, splitter)
}

func parseDiff(s string) (ret int, err error) {
	if s == "" {
		log.Println("hit")
	}
	fChar := s[0]
	if fChar == '+' {
		s = s[1:]
	}
	ret, err = strconv.Atoi(s)
	return
}

// This function splits a value from the original, and the value in the brackets, and returns both as integers
func splitValSubVal(s string) (first, bracketed int, err error) {
	first, err = strconv.Atoi(strings.Split(s, " ")[0])
	if err != nil {
		return
	}
	spl := strings.Split(s, "(")
	spl = strings.Split(spl[1], ")")
	bracketed, err = strconv.Atoi(spl[0])
	return
}

// ExtractMatch is used on a HLTV matchPage. It extracts all data from the matchpage, then will call functions to extract data from each player, and each map.
func ExtractMatch(url string) (match MatchData, err error) {

	// Initilizate MatchData with the data we have from just the URL.
	match = MatchData{
		MatchURL: url,
		MatchID:  strings.Split(url, "/")[4], //Extracting match id from URl, must have https:// prefixed
	}

	// Now we start scraping into the Struct
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0"))
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

		}
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
		match.Vetos = extractVetos(kids.Children().Text())
	})
	c.OnHTML(`.flexbox.left-right-padding`, func(e *colly.HTMLElement) {
		match.isDemo = true
		match.DemoLink = baseURL + e.Attr("href")
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
			pData := Player{
				TeamPlayedFor: curTeam,
			}
			extractPlayerData(e.ChildAttr(`a`, "href"), &pData)
			curTeam.Players = append(curTeam.Players, pData)
		})
		tCounter++
		return
	})

	// All player data from match page extracted. Now, for each map that was played, extract map data

	// Record time data was scraped
	match.ScrapedAt = time.Now()
	c.Visit(url)
	match.MapsPlayed = make([]MapData, len(match.MapLinks))
	for i, link := range match.MapLinks {
		match.MapsPlayed[i] = *ExtractStats(link)
	}
	return match, nil
}

func extractPlayerData(url string, p *Player) {
	p.PlayerURL = baseURL + url
	p.PlayerID = extractID(url)
	p.Name = strings.Split(url, "/")[3]
}

func extractVetos(data string) (result VetoList) {
	result = make(VetoList, 7)
	l := strings.Split(data, ".")
	if len(l) <= 1 && l[0] == "" {
		return
	}
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
				return
			}

		}
		var mname string
		if i == len(l)-1 {
			pb = 2
			mname = splitVeto[0]
		} else {
			mname = splitVeto[len(splitVeto)-1]
		}
		result[i-1] = veto{
			BanPick: pb,
			MapName: mname,
		}
	}
	return
}

func parseFloat(s string) (ret float32, err error) {
	prec, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return
	}
	// As we have chosen 32 bit , we can convert to a 32 bit float without fucking up the value.
	ret = float32(prec)
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

func parseErr(err error) error {

	if err != nil {
		return fmt.Errorf("error in extracting data from HTML: %s", err)
	}
	return nil

}
