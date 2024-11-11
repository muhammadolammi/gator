package main

import (
	"fmt"
	"log"
)

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Panicf("error getting next feed. err:%v", err)
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
		fmt.Println(item.Title)
	}
}
