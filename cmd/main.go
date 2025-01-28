package main

import (
	// "fmt"
	// "strings"
	"search_engine/web_scraper"
)

func main() {
	techSectionPageUrl := "https://www.bbc.com/amharic/topics/c06gq8wx467t?page="

	web_scraper.ScrapeArticleUrls(techSectionPageUrl, 22, "../corpus/tech_articles") // 22
	// fmt.Println(strings.Split("https://www.bbc.com/amharic/articles/cj0ry2jgzjpo", "/")[5])
}
