package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {

	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	// func (q *Queries) DeleteUsers(ctx context.Context) error {

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error deleting users: %w", err)
	}
	fmt.Printf("Deleted all users\n")
	return nil
}
