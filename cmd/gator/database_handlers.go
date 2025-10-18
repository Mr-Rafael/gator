package main

import (
	"fmt"
	"context"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error clearing the users table: %v", err)
	}
	err = s.db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error clearing the feeds table: %v", err)
	}
	err = s.db.ResetFeedFollows(context.Background())
	if err != nil {
		return fmt.Errorf("Error clearing the feed follows table: %v", err)
	}
	return nil
}