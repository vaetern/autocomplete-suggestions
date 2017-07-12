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

func findByLowestPrefixingDistance(suggestString string, trafficHubsList []trafficHub, trigramIndexList []trigramIndex) []trafficHub {

	//set lower\upper bounds for search by trigram index
	lowerBoundPosition, upperBoundPosition := getTrigramIndexes(suggestString, trigramIndexList)
	nextTrigramStartingOffset := upperBoundPosition

	var medianPosition int
	var needlePosition int8

	//binary search in bounds derived from trigram index
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

	//rewind to start of words cluster that satisfy search condition
	offsetToStartSuggestions := medianPosition
	for ; ; {
		needlePosition = whereIsNeedlePositionRelativeToString(trafficHubsList[medianPosition].name, suggestString)
		if needlePosition != 0 {
			offsetToStartSuggestions = medianPosition + 1
			break
		} else {
			medianPosition -= 1
		}
	}


	var tHub trafficHub

	//pick some(10, 15, HowManySuggestionsToReturn) traffic hubs that satisfies suggest conditions
	suggestedTrafficHubWithRangeList := []trafficHubWithRange{}
	bias := HowManySuggestionsToReturn
	for i := offsetToStartSuggestions; i < nextTrigramStartingOffset; i++ {
		tHub = trafficHubsList[i]
		//find distance in case insensitive manner
		stringRange := s.Index(s.ToLower(tHub.name), s.ToLower(suggestString))
		if stringRange >= 0 {
			suggestedTrafficHubWithRangeList = append(suggestedTrafficHubWithRangeList, trafficHubWithRange{tHub, stringRange})
		}
		if len(suggestedTrafficHubWithRangeList) > bias {
			break
		}
	}

	//sort suggested traffic hubs from lowest lexik distance to highest
	sort.Slice(suggestedTrafficHubWithRangeList, func(i, j int) bool {
		return suggestedTrafficHubWithRangeList[i].stringRange > suggestedTrafficHubWithRangeList[j].stringRange
	})

	//transform th+range in th
	suggestedTrafficHubList := trafficHubWithRangeListToTrafficHubList(suggestedTrafficHubWithRangeList)

	return suggestedTrafficHubList
}
