package hltvscrape

var ()

type SeriesScore int8
type MatchData struct { // @TODO add Map URLS into this struct. Change mapextractdata function to work from the statpage urls.
	// Data extracted from main match page.
	MatchURL         string      // URL of the match page.
	MatchID          string      // The ID of the match as located in the middle of 'url.../matches/{MATCHID}/north-vs...
	Team0            Team        // Team listed on the left side of HLTV
	Team1            Team        // Team listed on the right side of HLTV
	Team0SeriesScore SeriesScore // Map score for Team0
	Team1SeriesScore SeriesScore // Map Score for Team1
	MatchTimeEpoch   int         // Unix Time Match was played.
	Event            string      // What event the match was played in.
	EventID          string      // Id of the event located in the URL, similar to matchID

	BestOfType int    // Best of what? 3 or 1 or 5?
	Stage      string // What stage of the tournament the match was played in( semi final, final etc...)
	Winner     int8   // Team that won the game. 1  for Team0, 2 for Team1 and 0 for a draw.
	Vetos      VetoList
	MapLinks   []string
	MapsPlayed []MapData
	// Scrape Metadata
	ScrapedAtEpoc int // Unix Time webpage was scraped.
}

type MapData struct { // @TODO add selector strings like XML decoding. Example in: https://github.com/gocolly/colly/blob/master/_examples/hackernews_comments/hackernews_comments.go
	statPageURL string
	MapName     string
	Winner      int8
	Picker      int8

	Team0ScoreFirstHalf  int8
	Team0ScoreSecondHalf int8
	Team0ScoreTotal      int8
	Team0TeamRating      float32
	Team0FirstKills      int8
	Team0PlayerData      [4]PlayerMapPerf

	Team1ScoreFirstHalf  int8
	Team1ScoreSecondHalf int8
	Team1ScoreTotal      int8
	Team1TeamRating      float32
	Team1FirstKills      int8
	Team1PlayerData      [4]PlayerMapPerf
}

// VetoList is a map of the vetos of the match. Keyed by when each action was taken
type VetoList []veto

type veto struct {
	BanPick int8 // 0 if map picked, 1 if map banned, 2 if map left over
	MapName string
}

type PlayerMapPerf struct {
	Playername Player
	// Kills is the total kills INCLUDING headshots
	Kills     int8
	Headshots int8
	// Assists is the total assists INCLUDING flashassists
	Assists      int8
	FlashAssists int8
	Deaths       int8
	// KASTPercentage is the amount of rounds that the player got a Kill, Survived, an Assist or got Traded.
	KASTPercentage float32
	// KillDeathDiff is Kills - Deaths
	KillDeathDiff int8
	// ADR is the Average Damage per Round
	ADR float32
	// FirstKillsDiff is FirstKills - FirstDeaths
	FirstKillsDiff int8
	Rating         float32
}

// @TODO extract data from Team page.
type Team struct {
	TeamURL string
	TeamID  string
	Name    string
	Players []Player
}

type Player struct { //@TODO

}
