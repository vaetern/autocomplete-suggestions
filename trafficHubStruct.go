package main

type trafficHub struct {
	id   int
	name string
}

func trafficHubWithRangeListToTrafficHubList(thWithRangeList []trafficHubWithRange) []trafficHub {
	bridgedTrafficHubList := []trafficHub{}
	for i := 0; i < len(thWithRangeList); i++ {
		bridgedTrafficHubList = append(bridgedTrafficHubList, thWithRangeList[i].tHub)
	}

	return bridgedTrafficHubList
}
