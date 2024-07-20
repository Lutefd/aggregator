package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Lutefd/aggregator/internal/database"
	"github.com/google/uuid"
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
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}
		pubAt, err := parseDate(item.PubDate)
		if err != nil {
			log.Printf("error parsing time: %v\n", err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			PublishedAt: pubAt,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("error creating post: %s\n", err)
			continue
		}
	}
	if err != nil {
		log.Printf("error fetching feed %s: %s\n", feed.Url, err)
		return
	}
	log.Printf("fetched feed %s, %d ṕosts found", feed.Url, len(rssFeed.Channel.Item))
}
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC3339,
	}

	var t time.Time
	var err error

	for _, format := range formats {
		t, err = time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
