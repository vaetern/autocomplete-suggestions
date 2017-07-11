package main

import s "strings"

type trigramIndex struct {
	offset  int
	trigram string
}

func getTrigramIndexes(suggestString string, trigramIndexList []trigramIndex) (int, int) {
	indexLow, indexHigh := -1, -1

	for index, trigramIndex := range trigramIndexList {
		if s.EqualFold(trigramIndex.trigram, suggestString[0:3]) { //todo: out of bounds! panic maybe
			indexLow = trigramIndex.offset
			indexHigh = trigramIndexList[index+1].offset //todo: out of bounds!
			break
		}
	}

	return indexLow, indexHigh
}