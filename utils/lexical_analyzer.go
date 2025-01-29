package utils

import "strings"

var CleanTerms [][]TokenizedTerms

// takes in an array of terms (tokenized)
func DocumentLexicalAnalyzer() [][]TokenizedTerms {
	// removes stop words
	// replaces abbreviations
	// stems each terms
	articles := GetDocuments()
	for i, val := range articles {
		for j, val2 := range val {
			CleanTerms[i][j] = GetTerms(val2)
		}
	}

	return CleanTerms
}

func QueryLexicalAnalyzer(query string) {

}

func removeStopWords(term string) string {
	stopWords := GetStopWords()

	if _, found := stopWords[term]; found {
		return ""
	}
	return term
}

func replaceAbbreviations(term string) string {
	abbreviations := GetAbbreviations()

	if _, found := abbreviations[term]; found {
		return abbreviations[term]
	}
	return term
}

var unwanted_characters = []string{
	"“",
	"”",
	"።",
	":",
	"(",
	")",
}

func removeCharacters(term string) string {
	for _, val := range unwanted_characters {
		term = strings.ReplaceAll(term, val, "")
	}
	return term
}
