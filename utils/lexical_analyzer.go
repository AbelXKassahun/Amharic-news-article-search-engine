package utils

import "strings"

var CleanTerms [][]TokenizedTerms

// takes in an array of terms (tokenized)
func DocumentLexicalAnalyzer() [][]TokenizedTerms {
	// removes stop words
	// replaces abbreviations
	// stems each terms
	articles := GetDocuments()
	CleanTerms = make([][]TokenizedTerms, len(articles))
	for i, val := range articles {
		CleanTerms[i] = make([]TokenizedTerms, len(val))
		for j, val2 := range val {
			CleanTerms[i][j] = GetTerms(val2)
		}
	}

	return CleanTerms
}

func RemoveStopWords(term string) string {
	stopWords := GetStopWords()

	if _, found := stopWords[term]; found {
		return ""
	}
	return term
}

func ReplaceAbbreviations(term string) string {
	abbreviations := GetAbbreviations()

	if _, found := abbreviations[term]; found {
		return abbreviations[term]
	}
	return term
}

// do not include "/" or "-"
var unwanted_characters = []string{ 
	"“",
	"”",
	"\"",
	"።",
	":",
	"፡",
	"(",
	")",
}

func RemoveCharacters(term string) string {
	for _, val := range unwanted_characters {
		if strings.ContainsAny(term, val) {
			term = strings.ReplaceAll(term, val, "")
		}
	}
	return term
}
