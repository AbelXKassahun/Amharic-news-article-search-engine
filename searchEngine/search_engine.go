package searchEngine

import (
	// "fmt"
	"log"
	"math"
	"search_engine/utils"
	"sort"
	"strings"

	"github.com/AbelXKassahun/Amharic-Stemmer/stemmer"
)

type QueryFrequency struct {
	TF int
	Composite_Weight float64
}

type Query map[string]QueryFrequency

type QueryArticleSimilarity map[string]Similarity

type Similarity struct{
	Dot_Product float64
	Coefficient float64
}
// a type for sorting
type SimilarityEntry struct {
	ArticleID string
	Similarity Similarity
}

func SearchEngine(query string) []utils.Article{ // QueryArticleSimilarity
	// calls TermWeighing and gets the vector format of the query
	// calls TermWeighing and gets the vector format of every document
	queryTokens := tokenizeQuery(query)
	weightedQuery := computeQueryTermWeight(queryTokens)
	ranking := rankByCoefficient(removeZeroSimilarities(similarityComparison(weightedQuery)))

	// count := 0
	// fmt.Println("here ->")
	// for _, entry := range ranking {
	// 	if count  < 10 {
	// 		count++
	// 		fmt.Printf("ArticleID: %s, Coefficient: %v\n", entry.ArticleID, entry.Similarity.Coefficient) // %.2f
	// 	} else {
	// 		break
	// 	}
	// }
	// fmt.Println(count)

	return getFullArticleInfo(ranking[:10])
}

func getFullArticleInfo(ranking []SimilarityEntry) []utils.Article{
	documents := utils.GetDocuments()
	result := []utils.Article{}
	for _, val := range ranking {
		for _, articles := range documents {
			for _, article := range articles {
				if val.ArticleID == article.Article_ID {
					result = append(result, article) 
				}
			}
		}
	}
	return result
}

// Ranking function
func rankByCoefficient(similarities QueryArticleSimilarity) []SimilarityEntry {
	// Convert the map to a slice of SimilarityEntry
	entries := make([]SimilarityEntry, 0, len(similarities))
	for articleID, sim := range similarities {
		entries = append(entries, SimilarityEntry{
			ArticleID: articleID,
			Similarity: sim,
		})
	}

	// Sort the slice by Coefficient in descending order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Similarity.Coefficient > entries[j].Similarity.Coefficient
	})

	return entries
}

func removeZeroSimilarities(similarities QueryArticleSimilarity) QueryArticleSimilarity {
    filtered := make(QueryArticleSimilarity)

    for key, value := range similarities {
        // Check if both Dot_Product and Coefficient are non-zero
        if value.Dot_Product != 0 || value.Coefficient != 0 {
            filtered[key] = value
        }
    }

    return filtered
}

func tokenizeQuery(query string) Query{
	terms := strings.Split(query, " ")
	// for i := range terms{
	// 	fmt.Printf("[%s]\n",terms[i])
	// }
	// fmt.Println()
	query_frequency := make(Query)
	
	var newTerms []string
	for _, val := range terms {
		term := utils.ReplaceAbbreviations(utils.RemoveStopWords(utils.RemoveCharacters(val)))
		// fmt.Printf("[%s] -> [%v]\n", val, term)
		
		if term == "" {
			continue
		}
		splitted_term := strings.Split(term, " ")
		for j, val1 := range splitted_term {
			stemmed, err := stemmer.Stem(val1)
			if err != "" {
				log.Printf("Couldnt stem %s", val1)
				splitted_term[j] = val1	
			} else {
				splitted_term[j] = stemmed[0]
			}
		}
		newTerms = append(newTerms, splitted_term...)
	}
	terms = newTerms

	for _, val := range terms {
		if _, found := query_frequency[val]; found {
			freq := query_frequency[val]
			freq.TF++
			query_frequency[val] = freq
		} else {
			freq := query_frequency[val]
			freq.TF = 1
			query_frequency[val] = freq
		}
	}

	return query_frequency
}

func computeQueryTermWeight(queryTokens Query) Query{
	documentWeights := GetWeightForDocuments()
	for key, value := range queryTokens {
		for _, articles := range documentWeights {
			for _, terms := range articles {
				if _, found := terms[key]; found{
					weight := value
					weight.Composite_Weight = float64(weight.TF) * terms[key].IDF
					queryTokens[key] = weight
				}
			}
		}
	}
	return queryTokens
}

func computeQueryLength(queryTokens Query) float64 {
	var square float64
	for _, value := range queryTokens {
		if value.Composite_Weight != 0 {
			square += math.Pow(value.Composite_Weight, 2)
		}
	}
	return math.Sqrt(square)
}

func computeDocumentLength() map[string]float64 {
	documentWeights := GetWeightForDocuments()
	lengthOfArticles := make(map[string]float64)

	for _, articles := range documentWeights {
		for artcleID, terms := range articles {
			for _, value := range terms {
				lengthOfArticles[artcleID] += math.Pow(value.Composite_Weight, 2)
			}
		}
	}

	for key, value := range lengthOfArticles {
		lengthOfArticles[key] = math.Sqrt(value)
	}
	return lengthOfArticles
}

func similarityComparison(weightedQuery Query) QueryArticleSimilarity {
	documentWeights := GetWeightForDocuments()
	similarity := make(QueryArticleSimilarity)
	for key, value := range weightedQuery {
		for _, articles := range documentWeights {
			for artcleID, terms := range articles {
				if freq, found := terms[key]; found{
					dot_product := similarity[artcleID]
					dot_product.Dot_Product += value.Composite_Weight * freq.Composite_Weight
					similarity[artcleID] = dot_product
				}
			}
		}
	}
	
	queryLength := computeQueryLength(weightedQuery)
	documentsLength := computeDocumentLength()
	for key, value := range documentsLength {
		denom := queryLength * value
		llm := similarity[key]
		llm.Coefficient = llm.Dot_Product/denom
		similarity[key] = llm 
	}
	return similarity
}


