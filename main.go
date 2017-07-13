package main

import (
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"runtime"
	"unicode/utf8"
	"errors"
	"log"
)

const JaroWinklerTreshold = 0.8

const HowManySuggestionsToReturn = 10

const DataSourceName = "user:user@tcp(127.0.0.1:3306)/acr"

const MaxProcesses = 8

var trafficHubsList = []trafficHub{}

var trigramIndexList = []trigramIndex{}

func main() {

	runtime.GOMAXPROCS(MaxProcesses)

	hydrationService := hydrationService{"mysql", DataSourceName}

	trafficHubsList, trigramIndexList = hydrateDataFromDb(hydrationService)

	log.Println("Done hydrating pool")

	// --- debug ---
	//timeStart := time.Now()
	//result := findSuggestion("fffff", trafficHubsList, trigramIndexList)
	//
	//fmt.Println(result)
	//fmt.Println(time.Since(timeStart))
	// --- debug ---

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

	suggestString := r.URL.RawQuery

	printResult, err := getSuggestionForString(suggestString, trafficHubsList, trigramIndexList)

	if err != nil {
		fmt.Fprintln(w, "<b style=\"color: red;\">")
		fmt.Fprintln(w, err)
		fmt.Fprintln(w, "</b>")
	} else {
		fmt.Fprintln(w, printResult)
	}
}

func findSuggestion(suggestString string, trafficHubsList []trafficHub, trigramIndexList []trigramIndex) []trafficHub {
	result, err := findByLowestPrefixingDistance(suggestString, trafficHubsList, trigramIndexList)
	if len(result) > 0 && err == nil {
		return result
	}

	result = findIfJaroWinklerClose(suggestString, trafficHubsList)
	return result
}

func getSuggestionForString(suggestString string, trafficHubsList []trafficHub, trigramIndexList []trigramIndex) (string, error) {

	var printResult string
	var exception error

	if utf8.RuneCountInString(suggestString) >= 3 {

		start := time.Now()
		result := findSuggestion(suggestString, trafficHubsList, trigramIndexList)
		elapsed := time.Since(start)
		printResult = formatResult(result, elapsed)
		if elapsed.Seconds() > 0.01 {
			log.Println(elapsed.Seconds(), suggestString)
		}

	} else {

		exception = errors.New("Too short argument length")
		printResult = ""

	}

	return printResult, exception
}
