package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	rss := RSSFeed{}
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, http.NoBody)

	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Add("User-Agent", "gator")
	resp, err := client.Do(req)

	if err != nil {
		return &RSSFeed{}, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return &RSSFeed{}, err
	}
	// Unescape HTML entities in the Channel Title and Description
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	// Unescape HTML entities in each item
	for i := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
	}
	fmt.Println(rss)
	return &rss, nil
}
