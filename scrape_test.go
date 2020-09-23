// NEED TO EDIT INTO TEST FORMAT
func main() {

	c := colly.NewCollector(
		// Only visit hltv.org
		colly.AllowedDomains("www.hltv.org"),
		// Could look to implement parallelism as per below comment
		//colly.Async(true),
	)
	c.OnHTML(".matchstats", func(e colly.HTMLElement) {
		temp := item{}
		temp.MatchURL = 
		

}
