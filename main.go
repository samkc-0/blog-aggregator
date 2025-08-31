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
	cmds.Register("users", handleUsers)
	cmds.Register("reset", handlerReset)
	cmds.Register("agg", handlerAgg)
	cmds.Register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.Register("feeds", handlerFeeds)
	cmds.Register("follow", middlewareLoggedIn(handlerFollow))
	cmds.Register("following", middlewareLoggedIn(handlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.Register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("expected at least 1 command")
	}
	cmd := Command{Name: args[0], Args: args[1:]}
	if err := cmds.Run(&state, cmd); err != nil {
		log.Fatal(err)
	}
}
