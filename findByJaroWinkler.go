package main

import (
	"github.com/xrash/smetrics"
	"sort"
)

func findIfJaroWinklerClose(suggestString string, trafficHubsList []trafficHub) []trafficHub {

	rangeHubsList := []trafficHubWithRange{}
	for _, tHub := range trafficHubsList {
		stringRange := smetrics.JaroWinkler(suggestString, tHub.name, 0.9, 0)
		if stringRange > JaroWinklerTreshold {
			normalized := stringRange * 10000
			rangeHubsList = append(rangeHubsList, trafficHubWithRange{tHub, int(normalized)})
		}
	}

	sort.Slice(rangeHubsList, func(i, j int) bool {
		return rangeHubsList[i].stringRange > rangeHubsList[j].stringRange
	})

	resultHubsList := []trafficHub{}
	bias := 5
	for i := 0; i < len(rangeHubsList); i++ {
		resultHubsList = append(resultHubsList, rangeHubsList[i].tHub)
		if i == bias {
			break
		}
	}

	return resultHubsList

}