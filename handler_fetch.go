package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/whatsmynameagain/go-blog-aggregator/internal/database"
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

func handlerFetch(s *state, cmd command) error {
	hardcodedURL := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := fetchFeed(ctx, hardcodedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	// print output

	fmt.Printf("%v", feed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting URL: %w", err)
	}

	readResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	feed := RSSFeed{}

	err = xml.Unmarshal(readResp, &feed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

// untested
func scrapeFeeds(s *state) error {

	// nextFeed : database.Feed
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error trying to get next feed: %w", err)
	}

	fetchedParams := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: nextFeed.UpdatedAt,
	}

	err = s.db.MarkFeedFetched(context.Background(), fetchedParams)
	if err != nil {
		return fmt.Errorf("error marking feed fetched: %w", err)
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Println("Titles: ")
	for _, item := range feed.Channel.Item {
		fmt.Printf("%s\n", item.Title)
	}

	return nil
}
