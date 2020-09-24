package hltvscrape

import (
	"errors"

	"log"

	"github.com/gocolly/colly/v2"
)

// GetMapStatPage retreives map stats about a map that is played, the user specifies which map they want via an arg.
// It gets to the stats page from clicking through on the main match page, as HLTV stats page are retarded.
// Take matchURL, and which map we want to view.
func GetMapStatPage(matchURL string, mapNum int) (string, error) {
	// Declare a function to get the statpage from the mappage for the map we want.
	// Declare a variable to insert that URL into.
	c := colly.NewCollector() // No point in using caching here since it is very specific page we go to, which we only need to scrape/collect once.
	// There is a 'results-stats' a tag created for each map that is PLAYED.
	// We find both, and depending on value of 'mapNum' we visit that map.
	// If that map index does not exist, return that as an error.

	// This function will be called for each instance of the class it finds.
	// We use an interator to know how many elements we have passed, to ensure
	// that we are scraping the right map.
	elcount := 0
	var statURLS = make([]string, 7)
	// @TODO current issue where if mapNum = 0, we will goto 1, and then scrape that.
	c.OnHTML(`.results-stats`, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if elcount == mapNum {
			statURLS = append(statURLS, e.Request.AbsoluteURL(link))
		}
		elcount++
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting...")
	})
	c.Visit(matchURL)
	statPageURL := statURLS[mapNum]
	log.Printf("%v", statURLS)
	if statPageURL != "" {
		return statPageURL, nil
	}

	return "", errors.New("no Page found for that Map Number.\n")
}
