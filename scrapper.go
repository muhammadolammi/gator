package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/muhammadolammi/gator/internal/database"
)

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Println("error getting next feed. err: ", err)
		os.Exit(1)

	}
	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		log.Panicf("error marking feed as fetched. err:%v", err)
	}
	rssFeed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		log.Panicf("error fetching feeds for url(%v). err:%v", nextFeed.Url, err)
	}
	for _, item := range rssFeed.Channel.Item {
		postExist, err := s.db.PostExists(ctx, item.Link)
		if err != nil {
			log.Println("error checking if post exists err:", err)
		}
		if postExist {
			continue
		}
		var description sql.NullString
		if item.Description != "" {
			description = sql.NullString{
				Valid:  true,
				String: item.Description,
			}
		}
		publishedAt := time.Time{}
		formats := []string{
			time.RFC1123,          // "Mon, 02 Jan 2006 15:04:05 MST"
			time.RFC1123Z,         // "Mon, 02 Jan 2006 15:04:05 -0700"
			time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
			time.RFC850,           // "Monday, 02-Jan-06 15:04:05 MST"
			time.ANSIC,            // "Mon Jan _2 15:04:05 2006"
			time.UnixDate,         // "Mon Jan 2 15:04:05 MST 2006"
			"2006-01-02 15:04:05", // Common custom format "YYYY-MM-DD HH:MM:SS"
			"2006-01-02 15:04",    // "YYYY-MM-DD HH:MM"
			"02 Jan 2006 15:04",   // Custom format "DD Mon YYYY HH:MM"
			"02 Jan 06 15:04 MST", // Custom format with abbreviated year "DD Mon YY HH:MM TZ"
			"02 Jan 2006",         // "DD Mon YYYY"
		}
		parseSuccess := false
		for _, format := range formats {
			publishedAt, err = time.Parse(format, item.PubDate)
			if err == nil {
				parseSuccess = true
				break
			}
		}

		if !parseSuccess {
			log.Printf("error parsing published_at for item (%v): %v. Defaulting to current time", item.Title, err)
			publishedAt = time.Now()
		}

		post, err := s.db.CreatePost(ctx, database.CreatePostParams{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			FeedID:      nextFeed.ID,
			PublishedAt: publishedAt,
		})
		if err != nil {
			log.Println("error saving  post err:", err)
		}
		fmt.Println("post saved...")
		fmt.Println(post)

	}
}
