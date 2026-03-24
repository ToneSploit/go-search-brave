package main

import (
	"flag"
	"fmt"
	"log"

	"main/functions"
)

func main() {
	envFile := flag.String("env", "", "path to an alternative .env file (e.g. .env.nl)")
	flag.Parse()

	cfg, err := functions.LoadConfig(*envFile)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	results, err := functions.SearchCyberNews(cfg.BraveAPIKey, cfg.Keywords, &cfg.SearchOptions)
	if err != nil {
		log.Fatalf("search error: %v", err)
	}

	fmt.Printf("Fetched %d articles at %s\n", len(results.Articles), results.FetchedAt)
	for i, a := range results.Articles {
		fmt.Printf("Article %d:\n", i+1)
		fmt.Printf("  Title: %s\n", a.Title)
		fmt.Printf("    URL: %s\n", a.URL)

		if err := functions.SendPushoverNotification(&cfg.Pushover, &a); err != nil {
			log.Printf("pushover notification failed for article %d: %v", i+1, err)
		}
	}
}
