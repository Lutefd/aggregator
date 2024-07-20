package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Lutefd/aggregator/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		dbFeeds, err := db.GetNextFeedsToFecth(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error getting feeds to fetch: %s\n", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, dbFeed := range dbFeeds {
			wg.Add(1)
			go scrapeFeed(db, dbFeed, wg)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()
	err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error marking feed as fetched: %s\n", err)
		return
	}
	rssFeed, err := urlToRSSFeed(feed.Url)
	for _, item := range rssFeed.Channel.Item {
		log.Printf("item: %s\n", item.Title)
	}
	if err != nil {
		log.Printf("error fetching feed %s: %s\n", feed.Url, err)
		return
	}
	log.Printf("fetched feed %s\n", feed.Url)
}
