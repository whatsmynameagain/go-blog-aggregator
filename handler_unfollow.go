package main

import (
	"context"
	"fmt"

	"github.com/whatsmynameagain/go-blog-aggregator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	deleteParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteParams)
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %w", err)
	}

	return nil
}
