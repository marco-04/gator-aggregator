package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/marco-04/gator-aggregator/internal/database"
	"github.com/marco-04/gator-aggregator/internal/rss"
	"github.com/marco-04/gator-aggregator/internal/state"
)

func scrapeFeeds(s *state.State) error {
	nextFeed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	err = s.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fetchedFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("feed fetch error: %w", err)
	}

	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf("Fetched feed title: %s\n", item.Title)
	}

	return nil
}

func agg(s *state.State, args[]string) error {
	timeBetweenReqs := args[0]
	dur, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("duration string parsing error: %w", err)
	}

	fmt.Println("Collecting feeds every " + dur.String())
	ticker := time.NewTicker(dur)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}

