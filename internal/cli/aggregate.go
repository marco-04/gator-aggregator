package cli

import (
	"context"
	"fmt"

	"github.com/marco-04/gator-aggregator/internal/state"
	"github.com/marco-04/gator-aggregator/internal/rss"
)

func agg(s *state.State, args[]string) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(*feed)
	return nil
}

