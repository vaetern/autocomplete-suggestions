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
		objectId int
		content  string
	)
	rows, err := db.Query("select objectId, content from traffic_hub_translation WHERE content<>'' " +
	//"AND content LIKE 'Zar%'" +
		" ORDER BY content")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	trafficHubsList := []trafficHub{}
	trigramIndexList := []trigramIndex{}

	var slidingTrigram string

	log.Println("Start populating ", time.Since(timeStart))

	for rows.Next() {
		err := rows.Scan(&objectId, &content)
		if err != nil {
			log.Fatal(err)
		}
		//trigram index add block
		trafficHubsList = append(trafficHubsList, trafficHub{objectId, content})
		if len(content) > 2 && !s.EqualFold(slidingTrigram, content[0:3]) {
			slidingTrigram = s.ToLower(content[0:3])
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
