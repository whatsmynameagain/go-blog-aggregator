package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	err := s.conf.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("Username set to '%s'\n", cmd.Args[0])
	return nil
}
