package cli

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/marco-04/gator-aggregator/internal/database"
	"github.com/marco-04/gator-aggregator/internal/rss"
	"github.com/marco-04/gator-aggregator/internal/state"
)

// This is a cry for help
var timeLayoutStrings = []string {
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
}

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
		postUUID    := uuid.New()
		currentTime := time.Now()

		// Date strings were a mistake
		var parsedPubDate time.Time
		var err error
		for _, layout := range timeLayoutStrings {
			var date time.Time
			date, err = time.Parse(layout, item.PubDate)
			if err == nil {
				parsedPubDate = date
				break
			}
		}
		if err != nil {
			fmt.Printf("date parsing error: %v\n", err)
		}

		addedPost, err := s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID: postUUID,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Title: sql.NullString{
				String: item.Title,
				Valid: item.Title != "",
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid: item.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time: parsedPubDate,
				Valid: err == nil,
			},
			FeedID: nextFeed.ID,
		})

		if err != nil {
			// sheesh, there are no error enums for this library :|
			// If the error doesn't contain the 'url' column name
			// then it's _probably_ not a duplicate url error
			//
			// Error formatting seen at
			// `https://cs.opensource.google/go/go/+/refs/tags/go1.25.6:src/database/sql/sql.go;l=3365`
			if !strings.Contains(err.Error(), "violates unique constraint \"posts_url_key\"") {
				fmt.Println("db error: %w", err)
			} else {
				fmt.Printf("== Skipped %s ==\n", item.Link)
			}
		} else {
			fmt.Printf("Added new post: %v\n", addedPost)
		}
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

func browse(s *state.State, args[]string, user database.User) error {
	var limit int = 2
	var err error
	
	if len(args) >= 1 {
		limit, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("limit string parse error: %w", err)
		}
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Name: user.Name,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fmt.Printf("Showing %d posts:\n", limit)
	for _, post := range posts {
		fmt.Printf("* %+v\n", post)
	}

	return nil
}

