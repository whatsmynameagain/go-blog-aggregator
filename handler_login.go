package main

import (
	"context"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("could not find user: %w", err)
	}

	err = s.conf.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("Username set to '%s'\n", cmd.Args[0])
	return nil
}
