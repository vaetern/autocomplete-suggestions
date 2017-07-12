package main

import (
	s "strings"
	"unicode/utf8"
	"errors"
)

type trigramIndex struct {
	offset  int
	trigram string
}

func getTrigramIndexes(suggestString string, trigramIndexList []trigramIndex) (int, int, error) {
	indexLow, indexHigh := -1, -1
	var err error

	if utf8.RuneCountInString(suggestString) < 3{
		err = errors.New("Argument length too short")
		return indexLow, indexHigh, err
	}

	for index, trigramIndex := range trigramIndexList {
		if s.EqualFold(trigramIndex.trigram, suggestString[0:3]) { //todo: out of bounds! panic maybe
			indexLow = trigramIndex.offset
			indexHigh = trigramIndexList[index+1].offset //todo: out of bounds!
			break
		}
	}

	return indexLow, indexHigh, err
}