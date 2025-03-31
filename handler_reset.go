package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {

	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}
	fmt.Printf("Deleted all users\n")
	return nil
}
