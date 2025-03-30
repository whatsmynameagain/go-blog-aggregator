package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	funcs map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.funcs[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	funcToRun, ok := c.funcs[cmd.Name]
	if !ok {
		return fmt.Errorf("function for command '%s' does not exist", cmd.Name)
	}

	return funcToRun(s, cmd)
}
