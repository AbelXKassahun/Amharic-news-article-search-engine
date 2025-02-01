package utils

import (
	// "fmt"
	"log"
	"strings"
	"regexp"
	"github.com/AbelXKassahun/Amharic-Stemmer/stemmer"
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

	re := regexp.MustCompile(`^[a-zA-Z]+$`)

	var newTerms []string
	for _, val := range article.Content {
		if val.ContentType == "p" {
			sentence := strings.Split(val.Text, " ")
			for _, word := range sentence {
				term := ReplaceAbbreviations(RemoveStopWords(RemoveCharacters(word)))
				if term == "" {
					continue
				}
				splitted_term := strings.Split(term, " ")
				for j, val1 := range splitted_term {
					// fmt.Println(val1)
					if !re.MatchString(val1){
						stemmed, err := SafeCall(val1, stemmer.Stem)
						// fmt.Printf("stemmed[%v]\n", stemmed)
						// stemmed, err := stemmer.Stem(val1)
						if err != "" {
							log.Printf("Couldnt stem %s", val1)
							splitted_term[j] = val1	
						} else {
							splitted_term[j] = stemmed[0]
						}
					} else{
						splitted_term[j] = val1
					}
				}
				newTerms = append(newTerms, splitted_term...)
				// tokenizedTerms.Terms = append(tokenizedTerms.Terms, ReplaceAbbreviations(RemoveStopWords(RemoveCharacters(word))))
			}
		}
	}
	tokenizedTerms.Terms = newTerms
	return tokenizedTerms
}

func SafeCall(param string, fn func(string) ([]string, string)) (result []string, err string) {
	// Use defer to handle panics
	defer func() {
		if r := recover(); r != nil {
			// If a panic occurs, return the parameter itself
			// fmt.Println("Recovered from panic:", r)
			result = append(result, param)
			err = ""
		}
	}()

	// Call the function and return its result
	result, err = fn(param)
	return
}