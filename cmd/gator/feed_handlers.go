package main

import (
	"fmt"
	"context"
	"time"
	"strconv"
	"database/sql"
	"github.com/google/uuid"
	"github.com/Mr-Rafael/gator/internal/database"
	"github.com/Mr-Rafael/gator/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error: expected 1 argument (time between reqs), and found %v", len(cmd.Arguments))
	}
	duration, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("Error parsing the duration argument received: %v", err)
	}
	fmt.Printf("\nCollecting feeds every %v\n", duration)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for ;; <- ticker.C {
		fmt.Println("\nIt's scrapin' time!")
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("Error scraping feeds: %v", err)
		}
	}
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

func handlerBrowse(s *state, cmd command, userData database.User) error {
	var err error

	limit := 2
	if len(cmd.Arguments) >= 1 {
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("Expected a numeric parameter for command. %v", err)
		}
	}
	currentUser, err := getCurrentUserData(s)
	if err != nil {
		return fmt.Errorf("Error getting current user info: %v", err)
	}

	getPostsParams := database.GetPostsForUserParams{
		UserID: currentUser.ID,
		Limit: int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostsParams)
	if err != nil {
		return fmt.Errorf("Error getting posts for user: %v", err)
	}

	fmt.Println("Got the following posts:")
	for _, post := range posts {
		fmt.Printf(" - %v\n", post.Title)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	feedData, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed to update from the database: %v", err)
	}

	markFetchedParams := database.MarkFeedFetchedParams {
		ID: feedData.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err = s.db.MarkFeedFetched(context.Background(), markFetchedParams)
	if err != nil {
		return fmt.Errorf("Error marking the feed as updated in the database: %v", err)
	}	

	feedContent, err := rss.FetchFeed(context.Background(), feedData.Url)
	if err != nil {
		return fmt.Errorf("Error fetching the feed '%v' from URL: %v", feedData.Name, err)
	}

	fmt.Println("Storing the posts on the Database: ")

	for _, feedItem := range feedContent.Channel.Item {
		fmt.Printf("- Storing '%v'\n", feedItem.Title)
		savePostParams := database.CreatePostParams {
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: feedItem.Title,
			Url: feedItem.Link,
			Description: sql.NullString{
				String: feedItem.Description,
        		Valid: true,
    		},
			PublishedAt: parseNullableTime(feedItem.PubDate),
			FeedID:	feedData.ID,
		}
		create_error := s.db.CreatePost(context.Background(), savePostParams)
		if create_error != nil {
			return fmt.Errorf("Error storing the post on the database: %v", create_error)
		}
	}

	fmt.Println("\nDone scrapin'!\n")

	return nil
}

func parseNullableTime(input string) sql.NullTime {
	layouts := []string{
        time.RFC1123Z,
        time.RFC1123,
        time.RFC3339,
        "2006-01-02 15:04:05",
        "02 Jan 2006 15:04:05 MST",
        "Mon Jan 2 15:04:05 2006",
    }

    for _, layout := range layouts {
        if t, err := time.Parse(layout, input); err == nil {
            return sql.NullTime{Time: t, Valid: true}
        }
    }

    return sql.NullTime{Valid: false}

}
