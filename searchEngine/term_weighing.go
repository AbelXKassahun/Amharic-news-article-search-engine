package searchengine

import "sync"

type Vector map[string]float64

type Frequencies struct {
	frequency map[string][2]int
	mutex sync.Mutex
}

func GetWeightForDocuments() {
	// get weight of documents from json, using the open file funcion from the utility package

}

func TermWeighing(terms []string, isDocument bool) Vector {
	// compute tf * idf
	// needs N, dk, and tf. so it call getFrequencyOfATerm
	// saves the weight of each term in a file
	return Vector{}
}

var frequencies Frequencies // make sure its persitent and contains eveery unique term in every document

// gets the frequencies of a term
func getFrequencyOfATerm(term string, isDocument bool) {
	// check if a term's frequency (tf and dk) is already found
	// use mutex when reading and writing to the frequencies map
}