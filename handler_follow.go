package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/whatsmynameagain/go-blog-aggregator/internal/database"
)

/*
Add a follow command. It takes a single url argument and creates a new feed
follow record for the current user. It should print the name of the feed and
the current user once the record is created (which the query we just made should
support). You'll need a query to look up feeds by URL.
*/

/*
type GetFeedsRow struct {
	Name   string
	Url    string
	UserID uuid.UUID
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error)
*/

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to feed follow: %w", err)
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(ffRow.UserName, ffRow.FeedName)
	return nil
}

func handlerListFeedFollows(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
	if err != nil {
		return err
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
