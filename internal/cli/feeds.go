package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marco-04/gator-aggregator/internal/database"
	"github.com/marco-04/gator-aggregator/internal/state"
)

func addFeed(s *state.State, args[]string, user database.User) error {
	feedName    := args[0]
	feedURL     := args[1]
	feedUUID    := uuid.New()
	currentTime := time.Now()

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: feedUUID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name: feedName,
		Url: feedURL,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	err = followFeed(s, []string{ feedURL }, user)
	if err != nil {
		return fmt.Errorf("follow feed error: %w", err)
	}

	fmt.Printf("Feed [%s](%s) correctly added!\n", feedName, feedURL)
	fmt.Println(feed)

	return nil
}

func listFeeds(s *state.State, args[]string) error {
	feeds, err := s.DB.GetFeedsAndUser(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("* [%s](%s) (added by %s)\n", feed.Name, feed.Url, feed.UserName)
	}

	return nil
}

