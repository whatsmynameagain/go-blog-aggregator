package main

import (
	"context"
	"fmt"
)

func handlerUsers(s *state, cmd command) error {

	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	// func (q *Queries) GetUsers(ctx context.Context) ([]string, error) {

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting users list: %w", err)
	}

	for _, user := range users {
		fmt.Printf("%s", user)
		if user == s.conf.CurrentUserName {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}

	return nil
}
