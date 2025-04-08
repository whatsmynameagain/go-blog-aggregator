package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {

	// func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	// func (q *Queries) GetUserFromUUID(ctx context.Context, id uuid.UUID) (string, error) {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	ctx := context.Background()
	feedData, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds list: %w", err)
	}

	if len(feedData) == 0 {
		fmt.Printf("No feeds found\n")
		return nil
	}

	for _, feed := range feedData {
		userName, err := s.db.GetUserFromUUID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting username from UUID %v: %w", feedData[0].UserID, err)
		}
		fmt.Printf("%s: %s (%s) \n", feed.Name, feed.Url, userName)
	}

	return nil

}
