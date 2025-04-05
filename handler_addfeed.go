package main

import (
	"context"
	"fmt"

	"github.com/whatsmynameagain/go-blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	currentUserName := s.conf.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("error getting user %s: %w", currentUserName, err)
	}

	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: currentUser.ID,
	})

	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Printf("%v", newFeed)
	return nil
}
