package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"errors"
	"os"
	"context"
	"time"
	"database/sql"
	"github.com/google/uuid"
	"github.com/Mr-Rafael/gator/internal/config"
	"github.com/Mr-Rafael/gator/internal/database"
)

type state struct {
	db *database.Queries
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

	fmt.Println("Setting up valid commands")
	validCommands := make(map[string]func(*state, command) error)
	commands := commands{
		ValidCommands: validCommands,
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

	fmt.Println("Setting up state struct")
	fmt.Println("Reading configuration file")
	currentConf, err := config.Read()
	if err != nil {
		fmt.Printf("\nError reading configuration: %v", err)
	}
	fmt.Println("Opening the database")
	db, err := sql.Open("postgres", currentConf.DBURL)
	if err != nil {
		fmt.Printf("\nError connecting to the database: %v", err)
	}
	dbQueries := database.New(db)
	currentState := &state{
		Configuration: &currentConf,
		db: dbQueries,
	}
	fmt.Printf("Succesfully set up state.")

	fmt.Println("Reading user args")
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: Received less arguments than expected")
		os.Exit(1)
	}
	receivedCommand := getCommand(args)

	fmt.Println("Attempting to run the command")
	err =commands.run(currentState, receivedCommand)
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

	fmt.Printf("\nSearching for the user %v on the database.\n", userName)
	userData, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Error getting the user: %v", err)
	}
	fmt.Printf("\nFound the user: %s\n", userData)

	config.SetUser(userName)
	updateConfig(s)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("Error: expected an argument for the login function, and found 0")
	}
	
	creationTimeStamp := time.Now()
	fmt.Println("\nCreating the user at timestamp: %v\n", creationTimeStamp)
	userName := cmd.Arguments[0]
	creationParams := database.CreateUserParams {	
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: userName,
	}

	userData, err := s.db.CreateUser(context.Background(), creationParams)
	if err != nil {
		return fmt.Errorf("Error creating the user: %v", err)
	}

	config.SetUser(userName)
	updateConfig(s)
	fmt.Printf("The user was successfully created: %s", userData)
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