package main

import (
	"fmt"
	"time"
	"log"

	_ "github.com/go-sql-driver/mysql"
	//"os"
	"net/http"
	"runtime"
	//"math"
)

const JaroWinklerTreshold = 0.8

const HowManySuggestionsToReturn = 10

const DataSourceName = "user:user@tcp(127.0.0.1:3306)/acr"

var trafficHubsList = []trafficHub{}

var trigramIndexList = []trigramIndex{}

func main() {

	runtime.GOMAXPROCS(8)

	hydrationService := hydrationService{"mysql", DataSourceName}

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


	fmt.Fprintln(w, formatResult(result, elapsed))
}



func findSuggestion(suggestString string, trafficHubsList []trafficHub, trigramIndexList []trigramIndex) []trafficHub {
	result := findByLowestPrefixingDistance(suggestString, trafficHubsList, trigramIndexList)
	if len(result) > 0 {
		return result
	}

	result = findIfJaroWinklerClose(suggestString, trafficHubsList)
	return result
}