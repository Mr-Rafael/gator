package main

import (
	"fmt"
	"errors"
	"os"
	"github.com/Mr-Rafael/gator/internal/config"
)

type state struct {
	Configuration *config.Config
}

type command struct {
	Name string
	Arguments []string
}

type commands struct {
	ValidCommands map[string]func(*state, command) error
}

func main() {
	fmt.Println("Reading user args")
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: Received less arguments than expected")
		os.Exit(1)
	}
	loginCommand := getCommand(args)

	fmt.Println("Setting up valid commands")
	validCommands := make(map[string]func(*state, command) error)
	commands := commands{
		ValidCommands: validCommands,
	}
	commands.register("login", handlerLogin)

	fmt.Println("Setting up state struct")
	currentConf, err := config.Read()
	if err != nil {
		fmt.Printf("\nError reading configuration: %v", err)
	}
	currentState := &state{
		Configuration: &currentConf,
	}
	fmt.Printf("Succesfully set up state.")

	fmt.Println("Attempting to run command")
	err =commands.run(currentState, loginCommand)
	if err != nil {
		fmt.Printf("\nError running command: |%v|\n", err)
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("Error: expected an argument for the login function, and found 0")
	}
	userName := cmd.Arguments[0]

	config.SetUser(userName)
	updateConfig(s)
	return nil
}

func updateConfig(s *state) {
	updatedConfig, err := config.Read()
	if err != nil {
		fmt.Printf("\nError while updating config: %v", err)
	}
	s.Configuration = &updatedConfig
	fmt.Printf("\n|Successfully Updated Config|: %s\n", s.Configuration)
}

func getCommand(arguments []string) command {
	commandName := arguments[1]
	commandArguments := arguments[2:]
	return command{
		Name: commandName,
		Arguments: commandArguments,
	}
}

func (c *commands) run(s *state, cmd command) error {
	error := c.ValidCommands[cmd.Name](s, cmd)
	return error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.ValidCommands[name] = f
}