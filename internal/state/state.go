package state

import (
	"os"
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/marco-04/gator-aggregator/internal/config"
	"github.com/marco-04/gator-aggregator/internal/database"
)

var cfg config.Config

type State struct {
	DB     *database.Queries
	Config *config.Config
}

func Init() (State, error) {
	state := State {
		Config: &cfg,
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return State{}, fmt.Errorf("could not initialize state: %w", err)
	}

	dbQueries := database.New(db)
	state.DB = dbQueries

	return state, nil
}

func init() {
	err := config.InitIfNotExists()
	if err != nil {
		fmt.Printf("FATAL: could not write default config file: %v\n", err)
		os.Exit(1)
	}

	configFile, err := config.Read()
	if err != nil {
		fmt.Printf("FATAL: could not read config: %v\n", err)
		os.Exit(1)
	}
	cfg = configFile
}

