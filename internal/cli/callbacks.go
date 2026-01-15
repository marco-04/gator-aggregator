package cli

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/marco-04/gator-aggregator/internal/database"
	"github.com/marco-04/gator-aggregator/internal/state"
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
	"users": {
		description: "List users in the database",
		argNum: 0,
		usageStr: "",
		callback: listUsers,
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

func listUsers(s *state.State, args[]string) error {
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

