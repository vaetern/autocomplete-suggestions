package main

import (
	"fmt"
	"time"
	"log"

	"github.com/xrash/smetrics"

	_ "github.com/go-sql-driver/mysql"
	//"os"
	"sort"
	"net/http"
	"runtime"
	//"math"
)

const JaroWinklerTreshold = 0.8

var trafficHubsList = []trafficHub{}

var trigramIndexList = []trigramIndex{}

func main() {

	runtime.GOMAXPROCS(8)

	hydrationService := hydrationService{"mysql", "user:user@tcp(127.0.0.1:3306)/acr"}

	trafficHubsList, trigramIndexList = hydrateDataFromDb(hydrationService)

	//timeStart := time.Now()
	//result := findSuggestion("Paris+Airport", trafficHubsList, trigramIndexList)
	//
	//fmt.Println(result)
	//fmt.Println(time.Since(timeStart))

	http.HandleFunc("/", requestHandler) // each request calls requestHandler
	server := &http.Server{
		Addr:           ":8000",
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	suggestString := r.URL.RawQuery
	result := findSuggestion(suggestString, trafficHubsList, trigramIndexList)
	elapsed := time.Since(start)
	fmt.Fprintln(w, result)
	fmt.Fprintf(w, "Search process took %s", elapsed)
}



func findSuggestion(suggestString string, trafficHubsList []trafficHub, trigramIndexList []trigramIndex) []trafficHub {
	result := findByLowestPrefixingDistance(suggestString, trafficHubsList, trigramIndexList)
	if len(result) > 0 {
		return result
	}

	result = findIfJaroWinklerClose(suggestString, trafficHubsList)
	return result
}



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
