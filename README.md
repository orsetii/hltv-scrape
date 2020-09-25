# hltvscrape

A library that scrapes HLTV for various data relating to the CS:GO professional scene.

Written in go.

Download via `go get`:
```
go get -u github.com/orsetii/hltv-scrape
```

## @TODO 

- Finish data for past match page. 

- Create scraper for past matches, which then calls matchdata scrape on each one.

- Add functions to send data into CSV/other format(maybe)

- Figure out how to properly convert unix time without bugging out. Currently using epoch time as an int32.