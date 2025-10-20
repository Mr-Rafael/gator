package main

import (
	"fmt"
)

type command struct {
	Name string
	Arguments []string
}

type commands struct {
	ValidCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fmt.Printf("\nRunning command [%v] with parameters:\n", cmd.Name)
	for _, arg := range cmd.Arguments {
		fmt.Printf(" > %v\n", arg)
	}

	if handler, ok := c.ValidCommands[cmd.Name]; ok && handler != nil {
    	err := handler(s, cmd)
		if err != nil {
			fmt.Printf("\nError running the command: %v\n", err)
		}
	} else {
    	return fmt.Errorf("Unknown or unregistered command: '%v'", cmd.Name)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.ValidCommands[name] = f
}