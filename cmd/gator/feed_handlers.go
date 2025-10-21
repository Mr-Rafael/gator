package main

import (
	"fmt"
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/Mr-Rafael/gator/internal/database"
	"github.com/Mr-Rafael/gator/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("\nError fetching the feed: %v\n", err)
	}
	printStruct("Obtained the following feed:", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, userData database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("Error: expected 2 arguments (name, url), and found %v", len(cmd.Arguments))
	}
	feedName := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]

	feedCreationParams := database.CreateFeedParams {	
		ID: uuid.New(),
		Name: feedName,
		Url: feedURL,
		UserID: userData.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feedData, err := s.db.CreateFeed(context.Background(), feedCreationParams)
	if err != nil {
		return fmt.Errorf("Error inserting the feed: %v", err)
	}

	followCreationParams := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: userData.ID,
		FeedID: feedData.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), followCreationParams)
	if err != nil {
		return fmt.Errorf("Error creating follow in the database: %v", err)
	}

	printStruct("The feed was successfully created.", feedData)
	fmt.Printf("\nUser <%v> is now following '%v'.\n", userData.Name, feedData.Name)
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

func handlerFollow(s *state, cmd command, userData database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error: expected 1 argument (url), and found %v", len(cmd.Arguments))
	}
	feedURL := cmd.Arguments[0]

	feedData, err := s.db.GetFeedFromURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error getting the feed data: %v", err)
	}

	creationParams := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: userData.ID,
		FeedID: feedData.ID,
	}

	followData, err := s.db.CreateFeedFollow(context.Background(), creationParams)
	if err != nil {
		return fmt.Errorf("Error creating follow in the database: %v", err)
	}

	fmt.Printf("\nUser <%v> is now following '%v'.\n", followData.Name, followData.Name_2)
	return nil
}

func handlerFollowing(s *state, cmd command, userData database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), userData.ID)
	if err != nil {
		return fmt.Errorf("\nError fetching follow data: %v\n", err)
	}
	fmt.Printf("\nUser <%v> is following these feeds:\n", userData.Name)
	for _, follow := range feedFollows {
		fmt.Printf("\t- %v\n", follow.Name)
	}
	return nil
}


func handlerUnfollow(s *state, cmd command, userData database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error running 'following' command: 1 argument is needed, but received %v", len(cmd.Arguments))
	}
	feedURL := cmd.Arguments[0]

	feedData, err := s.db.GetFeedFromURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error getting the feed data: %v", err)
	}

	deleteParams := database.DeleteFollowParams {
		UserID: userData.ID,
		FeedID: feedData.ID,
	}
	err = s.db.DeleteFollow(context.Background(), deleteParams)
	if err != nil {
		return fmt.Errorf("\nError fetching follow data: %v\n", err)
	}
	fmt.Printf("\nUser <%v> is no longer following feed '%v'", userData.Name, feedData.Name)

	return nil
}