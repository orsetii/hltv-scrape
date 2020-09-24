package hltvscrape

import (
	"log"
	"testing"
)

func TestGetMap(t *testing.T) {
	// t.Fatal("not implemented")
	log.Println("Getting Map 0...")
	GetMapStatPage("https://www.hltv.org/matches/2344233/movistar-riders-vs-sj-dreamhack-open-fall-2020-closed-qualifier", 0)

	log.Println("Getting Map 1...")
	GetMapStatPage("https://www.hltv.org/matches/2344233/movistar-riders-vs-sj-dreamhack-open-fall-2020-closed-qualifier", 1)

}

func TestExtractMatch(t *testing.T) {
	log.Println("Extracting data from Match...")
	ret, _ := ExtractMatch("https://www.hltv.org/matches/2344233/movistar-riders-vs-sj-dreamhack-open-fall-2020-closed-qualifier")

	log.Printf("%v", ret)
}
