package main

import (
	"search_engine/web_scraper"
)

func main() {
	techSectionPageUrl := "https://www.bbc.com/amharic/topics/c06gq8wx467t?page="

	web_scraper.ScrapeArticleUrls(techSectionPageUrl, 1) // 22
}
