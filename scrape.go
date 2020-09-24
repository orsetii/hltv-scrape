package hltvscrape

import (
	"time"
)

var ()

type MatchData struct {
	// Data extracted from main match page.
	MatchURL   string    // URL of the match page.
	MatchID    string    // The ID of the match as located in the middle of 'url.../matches/{MATCHID}/north-vs...
	Team0      Team      // Team listed on the left side of HLTV
	Team1      Team      // Team listed on the right side of HLTV
	MatchTime  time.Time // Time Match was played.
	Event      string    // What event the match was played in.
	BestOfType int       // Best of what? 3 or 1 or 5?
	Winner     int8      // Team that won the game. 1  for Team0, 2 for Team1 and 0 for a draw.

	MapsPlayed []MapData
	// Scrape Metadata
	ScrapedAt time.Time // Time webpage was scraped.
}

type MapData struct {
	MapName string
	Winner  bool

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

type VetoList struct {
	Bans  map[int]MapBan  // Bans keyed by what stage of PB they were banned at.
	Picks map[int]MapPick // Picks keyed by what stage of PB they were picked at. This could be compressed further if not accounting for BO5 matches.
}
type MapBan struct {
	Team    bool
	MapName string
}
type MapPick struct {
	Team    bool
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

	Name string

	Players []Player
}

type Player struct { //@TODO
}
