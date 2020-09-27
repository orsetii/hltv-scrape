package hltvscrape

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

var mapIDs = map[string]int{
	"Mirage":      32,
	"Dust2":       31,
	"Nuke":        34,
	"Cobblestone": 39,
	"Overpass":    40,
	"Train":       35,
	"Cache":       29,
	"Inferno":     33,
	"Vertigo":     46,
}

func ExtractPastMatches(amount int, teamID string) (data []MatchData, err error) {
	data = make([]MatchData, amount)

	url := "https://www.hltv.org/results?team=" + teamID

	c := colly.NewCollector()

	collecIndex := 0
	c.OnHTML(`a.a-reset`, func(e *colly.HTMLElement) {
		if amount <= 0 {
			return
		}
		link := e.Attr("href")
		if !strings.Contains(link, "matches") {
			// If we get here we know it is not a correct link.
			return
		}
		data[collecIndex], err = ExtractMatch(BaseURL + link)
		if err != nil {
			log.Printf("Error in extracting match: %s", err)
		}
		collecIndex++
		amount--
	})
	c.Visit(url)
	return
}

func ExtractPastMaps(amount int, teamID string, teamName string, mapname string) (data []MapData, err error) {

	data = make([]MapData, amount)
	url := "https://www.hltv.org/stats/teams/map/" + strconv.Itoa(mapIDs[mapname]) + "/" + teamID + "/" + teamName

	c := colly.NewCollector()
	collecIndex := 0
	c.OnHTML(`td.time`, func(e *colly.HTMLElement) {
		val, ok := e.DOM.Children().Attr("href")
		if !ok {
			log.Printf("Unable to find link for map %d", collecIndex)
			return
		}
		data[collecIndex] = *ExtractStats(BaseURL + val)
		collecIndex++
	})
	c.Visit(url)
	return
}
