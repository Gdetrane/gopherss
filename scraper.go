package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Gdetrane/gopherss/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed '%s': %v", feed.Url, err)
	}

	for _, item := range rssFeed.Channel.Item {
    description := sql.NullString{}
    pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
    if err != nil {
      log.Printf("couldn't parse publishing date %v with err: %v", item.PubDate, err)
    }
    if item.Description != "" {
      description.String = item.Description
      description.Valid = true
    }
    _, dbErr := db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
      CreatedAt: time.Now().UTC(),
      UpdatedAt: time.Now().UTC(),
      Title: item.Title,
      Description: description,
      Url: item.Link,
      PublishedAt: pubAt,
      FeedID: feed.ID,
		})
    if dbErr != nil {
      if strings.Contains(dbErr.Error(), "duplicate key") {
        continue
      }
      log.Println("failed to create post:", err)
    }
	}
	log.Printf("Feed '%s' collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
