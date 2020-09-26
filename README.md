# hltvscrape

A library that scrapes HLTV for various data relating to the CS:GO professional scene.

Written in go.

Download via `go get`:
```
go get -u github.com/orsetii/hltv-scrape
```

## @TODO 

- Create scraper for past matches, which then calls matchdata scrape on each one. Needs to be able extract player data for the team past data we want to obtain data on. Create method on a map data to analyze the demo linked in the mapdata struct. This will call the defuselib parsing functions.

## Usage Notice

HLTV do not allow use of scraping without permission. Use at your own responbility.
