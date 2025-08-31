package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samkc-0/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, url string) (*RSSFeed, error) {
	httpClient := http.Client{Timeout: 10 * time.Second}

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "gator")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unsuccessful request with code %d fetching %s", response.StatusCode, url)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, err
	}

	rssFeed.UnescapeString()
	return &rssFeed, nil
}

func (f *RSSFeed) UnescapeString() {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)
	for i, item := range f.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		f.Channel.Item[i] = item
	}
}

func scrapeFeeds(s *State) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("failed to scrape feeds: %v", err)
	}
	if err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		log.Fatalf("failed to mark feed as fetched: %v", err)
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Printf("could not fetch %s from url '%s'\n", nextFeed.Name, nextFeed.Url)
		return
	}
	for _, item := range feed.Channel.Item {
		if err := s.savePostToDb(item, nextFeed.ID); err != nil {
			fmt.Printf("could not save item '%s' to db: %v\n", item.Title, err)
			fmt.Println("maybe it already exists")
			fmt.Println("---")
		}
	}
}

func (s *State) savePostToDb(r RSSItem, feedID uuid.UUID) error {
	pubDate, err := time.Parse(time.RFC3339, r.PubDate)
	pubDateAsNullTime := sql.NullTime{}
	if err != nil {
		pubDateAsNullTime.Time = pubDate
		pubDateAsNullTime.Valid = true
	}
	params := database.CreatePostParams{
		ID:          uuid.New(),
		Title:       r.Title,
		Url:         r.Link,
		Description: r.Description,
		PublishedAt: pubDateAsNullTime,
		FeedID:      feedID,
	}
	post, err := s.db.CreatePost(context.Background(), params)
	if err != nil {
		return fmt.Errorf("post '%s' (%s): %v", r.Title, r.Link, err)
	}
	fmt.Printf("saved post:\n%s\n%s\n---\n", post.Title, post.Url)
	return nil
}
