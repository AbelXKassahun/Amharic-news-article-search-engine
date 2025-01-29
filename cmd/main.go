package main

import (
	// "fmt"
	// "strings"
	// "search_engine/web_scraper"
	"search_engine/searchEngine"
	// "search_engine/utils"
)

func main() {
	// tech news pages url
	// techSectionPageUrl := "https://www.bbc.com/amharic/topics/c06gq8wx467t?page="
	// scrapes every article of every tech news page
	// web_scraper.ScrapeArticleUrls(techSectionPageUrl, 22, "../corpus/tech_articles") // 22
	// fmt.Println(utils.GetDocuments()[0][21].Article_ID)
	// utils.GetTerms(utils.GetDocuments()[0][21])
	searchengine.TermWeighing()
	searchengine.GetWeightForDocuments()
}
