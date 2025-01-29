package utils

import (
	"strings"
)

// "search_engine/web_scraper"

type TokenizedTerms struct {
	Article_ID string
	Terms      []string
}

func GetTerms(article Article) TokenizedTerms {
	tokenizedTerms := TokenizedTerms{}
	tokenizedTerms.Article_ID = article.Article_ID
	tokenizedTerms.Terms = make([]string, 0, 10)

	for _, val := range article.Content {
		if val.ContentType == "p" {
			sentence := strings.Split(val.Text, " ")
			for _, term := range sentence {
				tokenizedTerms.Terms = append(tokenizedTerms.Terms, replaceAbbreviations(removeStopWords(removeCharacters(term))))
			}
		}
	}
	return tokenizedTerms
}
