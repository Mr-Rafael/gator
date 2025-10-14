package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"errors"
	"os"
	"context"
	"time"
	"encoding/json"
	"database/sql"
	"github.com/google/uuid"
	"github.com/Mr-Rafael/gator/internal/config"
	"github.com/Mr-Rafael/gator/internal/database"
	"github.com/Mr-Rafael/gator/internal/rss"
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
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerFeeds)
	commands.register("reset", handlerReset)

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
	fmt.Println("Succesfully set up state")

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
	printStruct("Found the user:", userData)

	config.SetUser(userName)
	updateConfig(s)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("Error: expected an argument for the login function, and found 0")
	}
	
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
	printStruct("The user was successfully created:", userData)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	currentUser, err := config.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("Error reading the current user: %v", err)
	}

	usersData, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting the user list: %v", err)
	}

	for _, userData := range usersData {
		if userData.Name == currentUser {
			fmt.Printf("\n* %v (current)", userData.Name)
		} else {
			fmt.Printf("\n* %v", userData.Name)
		}
	}
	fmt.Println()
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error clearing the users table: %v", err)
	}
	err = s.db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error clearing the feeds table: %v", err)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("\nError fetching the feed: %v\n", err)
	}
	printStruct("Obtained the following feed:", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("Error: expected 2 arguments (name, url), and found %v", len(cmd.Arguments))
	}
	feedName := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]
	userName, err := config.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("Error reading the current user from configuration: %v", err)
	}
	userData, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Error querying the user data: %v", err)
	}

	creationParams := database.CreateFeedParams {	
		ID: uuid.New(),
		Name: feedName,
		Url: feedURL,
		UserID: userData.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feedData, err := s.db.CreateFeed(context.Background(), creationParams)
	if err != nil {
		return fmt.Errorf("Error inserting the feed: %v", err)
	}

	printStruct("The feed was successfully created", feedData)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feedsData, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to fetch data from the database: %v", err)
	}
	for _, feedData := range feedsData {
		printStruct("", feedData)
	}
	return nil
}

func updateConfig(s *state) {
	updatedConfig, err := config.Read()
	if err != nil {
		fmt.Printf("\nError while updating config: %v", err)
	}
	s.Configuration = &updatedConfig
	printStruct("Successfully updated config:", s.Configuration)
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
	if handler, ok := c.ValidCommands[cmd.Name]; ok && handler != nil {
    	handler(s, cmd)
	} else {
    	return fmt.Errorf("Unknown or unregistered command: '%v'", cmd.Name)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.ValidCommands[name] = f
}

func printStruct(description string, inter interface{}) {
	readable, err := json.MarshalIndent(inter, "", "  ")
	if err != nil {
		fmt.Printf("\nError preparing struct for printing: %v", err)
		return
	}
	fmt.Printf("\n%v\n%v\n", description, string(readable))
}