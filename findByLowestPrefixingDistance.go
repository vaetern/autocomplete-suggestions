package main

import ("sort"
		s "strings")


type trafficHubWithRange struct {
	tHub        trafficHub
	stringRange int
}

func whereIsNeedlePositionRelativeToString(subject string, needle string) int8 {
	longestStringLength := len(needle)
	if longestStringLength > len(subject) {
		longestStringLength = len(subject)
	}
	for i := 1; i <= longestStringLength; i++ {

		for !s.EqualFold(needle[i-1:i], subject[i-1:i]) {
			if s.ToLower(needle[i-1:i]) > s.ToLower(subject[i-1:i]) {
				return 1
			} else {
				return -1
			}
		}
	}
	return 0
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

func findByLowestPrefixingDistance(suggestString string, trafficHubsList []trafficHub, trigramIndexList []trigramIndex) []trafficHub {

	//set lower\upper bounds for search by trigram index
	lowerBoundPosition, upperBoundPosition := getTrigramIndexes(suggestString, trigramIndexList)
	upperSearchPosition := upperBoundPosition

	var medianPosition int
	var needlePosition int8

	for ; ; {
		medianPosition = (lowerBoundPosition + upperBoundPosition) / 2
		needlePosition = whereIsNeedlePositionRelativeToString(trafficHubsList[medianPosition].name, suggestString)
		if needlePosition == 1 {
			lowerBoundPosition = medianPosition + 1
		}
		if needlePosition == -1 {
			upperBoundPosition = medianPosition
		}
		if needlePosition == 0 || lowerBoundPosition == upperBoundPosition {
			break
		}
	}
	/////////
	offsetToStart := medianPosition
	for ; ; {
		needlePosition = whereIsNeedlePositionRelativeToString(trafficHubsList[medianPosition].name, suggestString)
		if needlePosition != 0 {
			offsetToStart = medianPosition + 1
			break
		} else {
			medianPosition -= 1
		}
	}

	var tHub trafficHub

	rangeHubsList := []trafficHubWithRange{}
	bias := 10
	for i := offsetToStart; i <= upperSearchPosition; i++ {
		tHub = trafficHubsList[i]
		stringRange := s.Index(tHub.name, suggestString)
		if stringRange >= 0 {
			rangeHubsList = append(rangeHubsList, trafficHubWithRange{tHub, stringRange})
		}
		if len(rangeHubsList) > bias {
			break
		}
	}

	sort.Slice(rangeHubsList, func(i, j int) bool {
		return rangeHubsList[i].stringRange > rangeHubsList[j].stringRange
	})

	resultHubsList := []trafficHub{}
	for i := 0; i < len(rangeHubsList); i++ {
		resultHubsList = append(resultHubsList, rangeHubsList[i].tHub)
	}

	return resultHubsList
}
