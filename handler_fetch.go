package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

	// check for arg (time_between_reqs)
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_requests>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	// moved to a different method to avoid unreachable code after the infinite loop
	runFetchLoop(s, timeBetweenReqs)

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
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    nextFeed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			// stealing this from the solution
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Printf("couldn't create post: %v", err)
			continue
		}

	}

	return nil
}

func runFetchLoop(s *state, timeBetweenReqs time.Duration) {
	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop() // not really needed because infinite loop
	for ; ; <-ticker.C {
		fmt.Printf("Scraping...\n") // temp, to make sure it's not bombarding the servers
		scrapeFeeds(s)
	}
}
