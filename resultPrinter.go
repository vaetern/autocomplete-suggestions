package main

import (
	"time"
)

func formatResult(resultList []trafficHub, elapsedTime time.Duration) string{

	resultingString := "- Search process took " + elapsedTime.String() + " -"

	for i := 0; i < len(resultList); i++ {
		resultingString += "\n" +
			resultList[i].name +
			"|"  +
			resultList[i].id
	}

	return resultingString
}
