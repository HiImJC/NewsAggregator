package main

import (
	"NewsAggregator/api"
	"NewsAggregator/pkg/aggregator"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

func main() {
	log.Println("Starting up aggregator. Golang version: ", runtime.Version())

	feedFile := flag.String("feedConfiguration", "./config/local.json", "path to feed configuration file")
	port := flag.Int("port", 8888, "port to start the application on")
	tick := flag.Duration("tick", 1*time.Hour, "duration between feed polls")
	flag.Parse()

	feeds, err := loadFeedConfiguration(*feedFile)
	if err != nil {
		log.Fatal("failed to load feed configuration " + err.Error())
	}

	a := aggregator.New(feeds, *tick)

	log.Fatal(api.StartServer(*port, a))
}

func loadFeedConfiguration(location string) (aggregator.Feeds, error) {
	c, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	var feeds aggregator.Feeds
	err = json.Unmarshal(c, &feeds)
	if err != nil {
		return nil, err
	}

	return feeds, nil
}
