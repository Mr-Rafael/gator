package main

import (
	"fmt"
	"errors"
	"context"
	"github.com/google/uuid"
	"time"
	"github.com/Mr-Rafael/gator/internal/config"
	"github.com/Mr-Rafael/gator/internal/database"
)

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