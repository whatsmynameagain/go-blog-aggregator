package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/whatsmynameagain/go-blog-aggregator/internal/config"
	"github.com/whatsmynameagain/go-blog-aggregator/internal/database"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	programState := &state{
		conf: &conf,
	}
	commands := commands{
		funcs: make(map[string]func(*state, command) error),
	}

	// reg login
	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	err = commands.run(programState, command{Name: commandName, Args: commandArgs})
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", programState.conf.DBURL)
	dbQueries := database.New(db)
	programState.db = dbQueries

}
