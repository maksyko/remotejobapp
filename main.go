package main

import (
	"github.com/ievgen-ma/remotejobapp/broadcasts"
	"github.com/ievgen-ma/remotejobapp/feeds"
	"github.com/ievgen-ma/remotejobapp/feeds/indeed"
	"github.com/ievgen-ma/remotejobapp/feeds/remoteglobal"
	"github.com/ievgen-ma/remotejobapp/feeds/remotive"
	"github.com/ievgen-ma/remotejobapp/feeds/stackoverflow"
	"github.com/ievgen-ma/remotejobapp/feeds/weworkremotely"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func createFeed(feedName string) feeds.PublicFeed {
	switch feedName {
	case "indeed":
		return indeed.NewPublicFeed(feedName)
	case "remoteglobal":
		return remoteglobal.NewPublicFeed(feedName)
	case "remotive":
		return remotive.NewPublicFeed(feedName)
	case "stackoverflow":
		return stackoverflow.NewPublicFeed(feedName)
	case "weworkremotely":
		return weworkremotely.NewPublicFeed(feedName)
	}
	return nil
}

func parseDate() {
	go createFeed("indeed").Connect()
	go createFeed("remoteglobal").Connect()
	go createFeed("remotive").Connect()
	go createFeed("stackoverflow").Connect()
	go createFeed("weworkremotely").Connect()
}

func broadcastData() {
	if err := broadcasts.NewBroadcastHandler().SendPosts(10); err != nil {
		log.Fatal(err)
	}

}

func main() {
	parseDate()

	go broadcastData()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
}
