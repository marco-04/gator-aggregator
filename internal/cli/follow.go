package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marco-04/gator-aggregator/internal/database"
	"github.com/marco-04/gator-aggregator/internal/state"
)

func followFeed(s *state.State, args[]string, user database.User) error {
	followUUID  := uuid.New()
	currentTime := time.Now()
	followURL   := args[0]

	feed, err := s.DB.GetFeedFromURL(context.Background(), followURL)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	feedFollow, err := s.DB.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID: followUUID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fmt.Printf("New feed subscription \"%s\" was successfully created for %s\n", followURL, s.Config.CurrentUserName)
	fmt.Println(feedFollow)

	return nil
}

func unfollowFeed(s *state.State, args[]string, user database.User) error {
	followURL   := args[0]

	feed, err := s.DB.GetFeedFromURL(context.Background(), followURL)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	err = s.DB.DeleteFollowFromUserAndFeed(context.Background(), database.DeleteFollowFromUserAndFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fmt.Printf("Feed subscription for \"%s\" was successfully deleted for %s\n", followURL, s.Config.CurrentUserName)

	return nil
}

func listFollowFeeds(s *state.State, args[]string) error {
	feeds, err := s.DB.GetFeedFollowsPerUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fmt.Printf("Feeds followed by %s:\n", s.Config.CurrentUserName)

	for _, feed := range feeds {
		fmt.Println("* " + feed)
	}

	return nil
}

