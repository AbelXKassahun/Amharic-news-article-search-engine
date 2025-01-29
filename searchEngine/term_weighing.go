package searchengine

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"search_engine/utils"
	"os"
)

type Documents map[string]map[string][4]float64

type TermFrequency map[int]map[string]map[string]Frequencies

type Frequencies struct {
	TF float64 
	DF int 
	IDF float64 
	Composite_Weight float64
}

var termFrequency = make(TermFrequency)

func TermWeighing()  {
	// compute tf * idf
	// needs N, dk, and tf. so it call getFrequencyOfATerm
	// saves the weight of each term in a file
	ComputeTermFrequency()
	ComputeWeight()

	printTermMeta(termFrequency)

	saveWeightToJSON()
}

func printTermMeta(meta TermFrequency) {
	for docIndex, articles := range meta {
		fmt.Printf("Document %d", docIndex)
		for articleID, termMap := range articles {
			fmt.Printf(" %s:\n",articleID)
			for term, freq := range termMap {
				fmt.Printf("    %s: {TF: %v, DF: %d, IDF: %v, weight: %v}\n", term, freq.TF, freq.DF, freq.IDF, freq.Composite_Weight)
			}
		}
	}
}


func ComputeTermFrequency(){
	// tokenized_terms := utils.DocumentLexicalAnalyzer()
	docs := [][]utils.TokenizedTerms{
		{// doc
			utils.TokenizedTerms{
				Article_ID: "d1a1",
				Terms: []string{
					"go", "is", "awesome", "go",
				},
			}, // article
			utils.TokenizedTerms{
				Article_ID: "d1a2",
				Terms: []string{
					"learning", "go", "is", "fun",
				},
			},
		},
		{
			utils.TokenizedTerms{
				Article_ID: "d2a1",
				Terms: []string{
					"go", "is", "powerful",
				},
			}, // article
			utils.TokenizedTerms{
				Article_ID: "d2a2",
				Terms: []string{
					"code", "in", "go", "daily",
				},
			},
		},
	}
	// an array of docs [][]TokenizedTerms
	// docs - an array of articles []TokenizedTerms
	// article - TokenizedTerms
	dfTempTracker := make(map[string]map[string]bool)
	for docIndex, articles := range docs { // tokenized_terms
		termFrequency[docIndex] = make(map[string]map[string]Frequencies)
		for _, article := range articles {
			termFrequency[docIndex][article.Article_ID] = make(map[string]Frequencies)
			for _, term := range article.Terms {
				freq := termFrequency[docIndex][article.Article_ID][term] 
				freq.TF++
				termFrequency[docIndex][article.Article_ID][term] = freq

				if dfTempTracker[term] == nil {
					dfTempTracker[term] = make(map[string]bool)
				}
				dfTempTracker[term][article.Article_ID] = true
			}
		}
	}

	for term, articlesContainingTerm := range dfTempTracker {
		for docIndex, articles := range termFrequency {
			for articleId, terms := range articles {
				if freq, found := terms[term]; found {
					freq.DF = len(articlesContainingTerm)
					termFrequency[docIndex][articleId][term] = freq
				}
			}
		}
	}
}


// compute the composite weight (tf * idf)
func ComputeWeight() {
	for docIndex, articles := range termFrequency {
		for articleID, terms := range articles {
			for term, freq := range terms {
				freq.TF = freq.TF/float64(len(terms)) // length normalized frequency
				if freq.DF > 0 {
					freq.IDF = math.Log2(float64(450/freq.DF))
					freq.Composite_Weight = float64(freq.TF) * freq.IDF
				} else {
					log.Fatalln("Found a document frequency with value less than or equal to 0")
				}
				
				termFrequency[docIndex][articleID][term] = freq
			}
		}
	}
}

func saveWeightToJSON(){
	updatedJSON, err := json.MarshalIndent(termFrequency, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding termFreqeuncy to JSON: %v\n", err)
	}
	if err := os.WriteFile("../utils/weight.json", updatedJSON, 0644); err != nil {
		log.Fatalf("Error writing file to ../utils/weight.json: %v\n", err)
	}
	log.Println("File created and data written to: ", "../utils/weight.json")
}

func GetWeightForDocuments() {
	gg := make(TermFrequency)

	// get weight of documents from json, using the open file funcion from the utility package
	file, err := os.ReadFile("../utils/weight.json")
	if err != nil {
		log.Fatal("Error reading weight.json: ", err)
	}

	// Unmarshal existing data into a slice
	if err := json.Unmarshal(file, &gg); err != nil {
		log.Fatal("Error decoding weight.json: ", err)
	}

	fmt.Println("Reading from file")
	printTermMeta(gg)
}
