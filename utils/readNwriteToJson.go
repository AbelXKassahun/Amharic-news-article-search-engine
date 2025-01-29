package utils

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"search_engine/web_scraper"
)

//go:embed stop-words.json
var stop_words []byte

//go:embed abbreviations.json
var abbreviations []byte

type StopWords struct {
	StopWords []string `json:"stop_words"`
}

type stop_word_map map[string]struct{}

type Abbreviaitons map[string]string

type Article web_scraper.Article

func GetStopWords() stop_word_map {
	var data StopWords
	err := json.Unmarshal(stop_words, &data)
	if err != nil {
		log.Fatalln("Error parsing JSON:", err)
	}

	stopWordsMap := make(stop_word_map)
	for _, word := range data.StopWords {
		stopWordsMap[word] = struct{}{}
	}
	return stopWordsMap
	// Example lookup
	// wordToCheck := "ስለሚሆን"
	// if _, found := stopWordsMap[wordToCheck]; found {
	// 	log.Printf("'%s' is a stop word.\n", wordToCheck)
	// } else {
	// 	log.Printf("'%s' is not a stop word.\n", wordToCheck)
	// }
}

func GetAbbreviations() Abbreviaitons {
	var data Abbreviaitons
	err := json.Unmarshal(abbreviations, &data)
	if err != nil {
		log.Fatalln("Error parsing JSON:", err)
	}
	return data

	// Example lookup
	// key := "ደ/ዘይት"
	// if value, found := abbreviations[key]; found {
	// 	fmt.Printf("'%s' means '%s'.\n", key, value)
	// } else {
	// 	fmt.Printf("'%s' not found in abbreviations.\n", key)
	// }
}

func GetDocuments() [][]Article{
	articles := [][]Article{}
	for i := 0; i < 1; i++ { // i < 22
		articles = append(articles, GetPage("tech_articles", fmt.Sprintf("page_%d.json", i+1)))
	}
	return articles
}

func GetPage(docname, filename string) []Article {
	articles := []Article{}
	filePath := filepath.Join("../corpus/", docname, filename)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error opening document %s: %v\n", filename, err)
	}

	if err := json.Unmarshal(file, &articles); err != nil {
		log.Fatal("Error decoding JSON: ", err)
	}
	return articles
}
