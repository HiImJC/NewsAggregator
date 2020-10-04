package main

import (
	"NewsAggregator/api"
	"NewsAggregator/pkg/aggregator"
	"log"
	"runtime"
	"time"
)

func main() {
	log.Println("Starting up aggregator. Golang version: ", runtime.Version())

	a := aggregator.New([]string{
		"http://feeds.bbci.co.uk/news/uk/rss.xml",
		"http://feeds.bbci.co.uk/news/technology/rss.xml",
		"http://feeds.skynews.com/feeds/rss/uk.xml",
		"http://feeds.skynews.com/feeds/rss/technology.xml",
	}, 5*time.Second)

	log.Fatal(api.StartServer(8080, a))
}
