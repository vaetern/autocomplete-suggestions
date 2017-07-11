package main

import (
	"time"
	"strconv"
)

func formatResult(resultList []trafficHub, elapsedTime time.Duration) string{

	resultingString := "- Search process took " + elapsedTime.String() + " -"

	for i := 0; i < len(resultList); i++ {
		resultingString += "\n" +
			resultList[i].name +
			"|"  +
			strconv.Itoa(resultList[i].id)
	}

	return resultingString
}
