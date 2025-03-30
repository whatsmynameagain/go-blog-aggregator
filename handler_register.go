package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/whatsmynameagain/go-blog-aggregator/internal/database"
)

func handlerRegister(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}
	err = s.conf.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}
	fmt.Printf("User created: '%v'\n", user)
	return nil
}
