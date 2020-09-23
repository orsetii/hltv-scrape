package main

import (
	"github.com/gocolly/colly/v2"
	"time"
)

type MatchData struct {
	MatchID   string    // The ID of the match as located in the middle of 'url.../matches/{MATCHID}/north-vs...
	Team1     string    // Team listed on the left side of HLTV
	Team2     string    // Team listed on the right side of HLTV
	MatchTime time.Time // Time Match was played.

	ScrapedAt time.Time // Time webpage was scraped.
}

type MapData struct {
	MapName           string
	TScoreFirstHalf   uint8
	TScoreSecondHalf  uint8
	CTScoreFirstHalf  uint8
	CTScoreSecondHalf uint8
	TScoreTotal       uint8
	CTScoreTotal      uint8
}
