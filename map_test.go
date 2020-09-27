package hltvscrape

import (
	"log"
	"testing"
)

func TestExtractMatch(t *testing.T) {
	// log.Println("Extracting data from Match...")
	// ret, err := ExtractPastMatches(3, "4579")
	// log.Printf("Error in getting past matches: %s", err)
	// for _, v := range ret {
	// 	logData(v)
	// }
	//ret, _ := ExtractMatch("https://www.hltv.org/matches/2344119/beyond-vs-checkmate-perfect-world-asia-league-fall-2020")
	//logData(ret)

	ret, _ := ExtractPastMaps(5, "4579", "LowLandLions", "Dust2")
	log.Printf("%+v", ret)
}

func logData(ret MatchData) {
	log.Printf("%+#v", ret)
	log.Printf("Loaded Match %s vs %s", ret.Team0.Name, ret.Team1.Name)
	log.Printf("Winner: Team%d", ret.Winner-1)
	log.Printf("Team0 Score: %d", ret.Team0SeriesScore)
	log.Printf("Team1 Score: %d", ret.Team1SeriesScore)
	log.Println("----DATA PROCESSING FINISHED----")
}
