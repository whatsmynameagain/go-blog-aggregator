package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(newFeed, currentUser)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
