package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/samkc-0/gator/internal/config"
	"github.com/samkc-0/gator/internal/database"
)

type State struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	queries := database.New(db)
	state := State{cfg: &cfg, db: queries}

	cmds := Commands{
		registered: make(map[string]func(*State, Command) error),
	}
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("expected at least 2 arguments")
	}
	cmd := Command{Name: args[0], Args: args[1:]}
	if err := cmds.Run(&state, cmd); err != nil {
		log.Fatal(err)
	}
}
