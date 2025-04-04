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

	db, err := sql.Open("postgres", conf.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		db:   dbQueries,
		conf: &conf,
	}
	commands := commands{
		funcs: make(map[string]func(*state, command) error),
	}

	// reg login
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerFetch)

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

}
