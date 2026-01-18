package cli

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/marco-04/gator-aggregator/internal/database"
	"github.com/marco-04/gator-aggregator/internal/state"
	"github.com/marco-04/gator-aggregator/internal/rss"
)

var cmds = map[string]cmd {
	"login": {
		description: "Set current user",
		argNum: 1,
		usageStr: "<user>",
		callback: login,
	},
	"dburl": {
		description: "Set db_url in config file",
		argNum: 1,
		usageStr: "<db_url>",
		callback: setDBURL,
	},
	"register": {
		description: "Register a user in the database",
		argNum: 1,
		usageStr: "<user>",
		callback: register,
	},
	"reset": {
		description: "[DEBUG] Reset database to allow for easier testing",
		argNum: 0,
		usageStr: "",
		callback: reset,
	},
	"users": {
		description: "List users in the database",
		argNum: 0,
		usageStr: "",
		callback: listUsers,
	},
	"agg": {
		description: "Fetch feeds",
		argNum: 0,
		usageStr: "",
		callback: agg,
	},
	"addfeed": {
		description: "Add a feed",
		argNum: 2,
		usageStr: "<name> <url>",
		callback: addFeed,
	},
	"feeds": {
		description: "List all feeds for the current user",
		argNum: 0,
		usageStr: "",
		callback: listFeeds,
	},
	"follow": {
		description: "Follow a feed",
		argNum: 1,
		usageStr: "<url>",
		callback: followFeed,
	},
	"following": {
		description: "List all feeds the current user is following",
		argNum: 0,
		usageStr: "",
		callback: listFollowFeeds,
	},
	"help": {
		description: "Print usage text",
		argNum: 0,
		usageStr: "",
		callback: printHelp,
	},
}

func login(s *state.State, args []string) error {
	currentUser := args[0]

	names, err := s.DB.GetUserNames(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	if !slices.Contains(names, currentUser) {
		return fmt.Errorf("User \"%s\" does not exist", currentUser)
	}

	if err := s.Config.SetUser(currentUser); err != nil {
		return err
	}

	fmt.Printf("\"%s\" set as current_user_name\n", currentUser)
	return nil
}

func setDBURL(s *state.State, args []string) error {
	dbURL := args[0]

	if err := s.Config.SetDBURL(dbURL); err != nil {
		return err
	}

	fmt.Printf("\"%s\" set as db_url\n", dbURL)
	return nil
}

func register(s *state.State, args []string) error {
	name := args[0]
	id := uuid.New()

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID: id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	if err := s.Config.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("\"%s\" successfully created!\n", name)
	fmt.Printf("User data: %v\n", user)
	return nil
}

func reset(s *state.State, args[]string) error {
	err := s.DB.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fmt.Println("Successfully cleared users table")
	return nil
}

func listUsers(s *state.State, args[]string) error {
	users, err := s.DB.GetUserNames(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("* %s", user)
		if user == s.Config.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func agg(s *state.State, args[]string) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(*feed)
	return nil
}

func addFeed(s *state.State, args[]string) error {
	feedName    := args[0]
	feedURL     := args[1]
	feedUUID    := uuid.New()
	currentTime := time.Now()

	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	var userUUID uuid.UUID
	for _, user := range users {
		if s.Config.CurrentUserName == user.Name {
			userUUID = user.ID
			break
		}
	}
	if userUUID == [16]byte{} {
		return fmt.Errorf("user %s was not found", s.Config.CurrentUserName)
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: feedUUID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name: feedName,
		Url: feedURL,
		UserID: userUUID,
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	err = followFeed(s, []string{ feedURL })
	if err != nil {
		return fmt.Errorf("follow feed error: %w", err)
	}

	fmt.Printf("Feed [%s](%s) correctly added!", feedName, feedURL)
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

func followFeed(s *state.State, args[]string) error {
	followUUID  := uuid.New()
	currentTime := time.Now()
	followURL   := args[0]

	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	var userUUID uuid.UUID
	for _, user := range users {
		if s.Config.CurrentUserName == user.Name {
			userUUID = user.ID
			break
		}
	}
	if userUUID == [16]byte{} {
		return fmt.Errorf("user %s was not found", s.Config.CurrentUserName)
	}

	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	var feedUUID uuid.UUID
	for _, feed := range feeds {
		if followURL == feed.Url {
			feedUUID = feed.ID
			break
		}
	}
	if userUUID == [16]byte{} {
		return fmt.Errorf("user %s was not found", s.Config.CurrentUserName)
	}

	feedFollow, err := s.DB.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID: followUUID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID: userUUID,
		FeedID: feedUUID,
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	fmt.Printf("New feed subscription \"%s\" was successfully created for %s", followURL, s.Config.CurrentUserName)
	fmt.Println(feedFollow)

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

