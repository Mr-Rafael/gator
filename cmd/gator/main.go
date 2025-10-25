package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"os"
	"context"
	"encoding/json"
	"database/sql"
	"github.com/Mr-Rafael/gator/internal/config"
	"github.com/Mr-Rafael/gator/internal/database"
)

type state struct {
	db *database.Queries
	Configuration *config.Config
}

func main() {
	validCommands := make(map[string]func(*state, command) error)
	commands := commands{
		ValidCommands: validCommands,
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))
	commands.register("reset", handlerReset)

	currentConf, err := config.Read()
	if err != nil {
		fmt.Printf("\nError reading configuration: %v", err)
	}
	db, err := sql.Open("postgres", currentConf.DBURL)
	if err != nil {
		fmt.Printf("\nError connecting to the database: %v", err)
	}
	dbQueries := database.New(db)
	currentState := &state{
		Configuration: &currentConf,
		db: dbQueries,
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: Received less arguments than expected")
		os.Exit(1)
	}
	receivedCommand := getCommand(args)

	err =commands.run(currentState, receivedCommand)
	if err != nil {
		fmt.Printf("\nError running command: |%v|\n", err)
		os.Exit(1)
	}
}

func getCurrentUserData(s * state) (database.User, error) {
	userName, err := config.GetCurrentUser()
	if err != nil {
		return database.User{}, fmt.Errorf("Error reading the current user from configuration: %v", err)
	}
	userData, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return database.User{}, fmt.Errorf("Error querying the user data: %v", err)
	}
	return userData, nil;
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

func printStruct(description string, inter interface{}) {
	readable, err := json.MarshalIndent(inter, "", "  ")
	if err != nil {
		fmt.Printf("\nError preparing struct for printing: %v", err)
		return
	}
	fmt.Printf("\n%v\n%v\n", description, string(readable))
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	currentConf, err := config.Read()
	if err != nil {
		fmt.Printf("\nError reading configuration: %v", err)
		return nil
	}
	db, err := sql.Open("postgres", currentConf.DBURL)
	if err != nil {
		fmt.Printf("\nError connecting to the database: %v", err)
		return nil
	}
	dbQueries := database.New(db)
	
	userName, err := config.GetCurrentUser()
	if err != nil {
		return nil
	}
	userData, err := dbQueries.GetUser(context.Background(), userName)
	if err != nil {
		return nil
	}

	return func(s *state, cmd command) error {
		return handler(s, cmd, userData)
	}
}