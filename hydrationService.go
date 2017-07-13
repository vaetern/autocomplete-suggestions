package main

import (
	"time"
	"database/sql"
	"log"
	s "strings"
)

type hydrationService struct {
	driverName     string
	dataSourceName string
}

func hydrateDataFromDb(service hydrationService) ([]trafficHub, []trigramIndex) {

	timeStart := time.Now()

	db, err := sql.Open(service.driverName, service.dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	var (
		uic  string
		name string
	)
	rows, err := db.Query("SELECT `uic`, `name` FROM traffic_hub WHERE `name`<>'' AND `uic` IS NOT NULL " +
	//"AND name LIKE 'Zar%'" +
		" ORDER BY name")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	trafficHubsList := []trafficHub{}
	trigramIndexList := []trigramIndex{}

	var slidingTrigram string

	log.Println("Start populating ", time.Since(timeStart))

	for rows.Next() {
		err := rows.Scan(&uic, &name)
		if err != nil {
			log.Fatal(err)
		}
		//trigram index add block
		trafficHubsList = append(trafficHubsList, trafficHub{uic, name})
		if len(name) > 2 && !s.EqualFold(slidingTrigram, name[0:3]) {
			slidingTrigram = s.ToLower(name[0:3])
			trigramIndexList = append(trigramIndexList, trigramIndex{len(trafficHubsList) - 1, slidingTrigram})
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Ready in ", time.Since(timeStart))
	return trafficHubsList, trigramIndexList
}
