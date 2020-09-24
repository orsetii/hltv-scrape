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
	logData(ret)
	ret, _ = ExtractMatch("https://www.hltv.org/matches/2344232/x6tence-vs-mousesports-dreamhack-open-fall-2020-closed-qualifier")
	logData(ret)
	ret, _ = ExtractMatch("https://www.hltv.org/matches/2344119/beyond-vs-checkmate-perfect-world-asia-league-fall-2020")
	logData(ret)
}

func logData(ret MatchData) {

	log.Printf("%+v", ret)
	log.Printf("Loaded Match %s vs %s", ret.Team0.Name, ret.Team1.Name)
	log.Printf("Winner: %v", ret.Winner)
	log.Printf("Team0 Score: %d", ret.Team0SeriesScore)
	log.Printf("Team1 Score: %d", ret.Team1SeriesScore)
}
