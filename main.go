package main

import (
	"log"
	"os"

	"github.com/samkc-0/gator/internal/config"
)

type State struct {
	cfg *config.Config
}

func main() {
	cfg := config.Read()
	state := State{cfg: &cfg}
	cmds := Commands{
		registered: make(map[string]func(*State, Command) error),
	}
	cmds.Register("login", handlerLogin)
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("expected at least 2 arguments")
	}
	cmd := Command{Name: args[0], Args: args[1:]}
	if err := cmds.Run(&state, cmd); err != nil {
		log.Fatal(err)
	}
}
