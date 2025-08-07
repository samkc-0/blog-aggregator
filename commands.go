package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/samkc-0/gator/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	registered map[string]func(*State, Command) error
}

func (c *Commands) Run(state *State, cmd Command) error {
	handler, ok := c.registered[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command %s", cmd.Name)
	}
	if err := handler(state, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	c.registered[name] = handler
}

func handlerLogin(state *State, cmd Command) error {
	if cmd.Name != "login" {
		return fmt.Errorf("expected login command, got %s", cmd.Name)
	}
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login command expected 1 argument (username), got %d arguments", len(cmd.Args))
	}
	username := cmd.Args[0]
	user, err := state.db.GetUser(context.Background(), username)
	if err != nil {
		log.Fatal(err)
	}
	if err := state.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Printf("logged in as user %s\n", state.cfg.CurrentUsername)
	return nil
}

func handlerRegister(state *State, cmd Command) error {
	if cmd.Name != "register" {
		return fmt.Errorf("expected register command, got %s", cmd.Name)
	}
	if len(cmd.Args) != 1 {
		return fmt.Errorf("register command expected 1 argument (username), got %d arguments", len(cmd.Args))
	}
	username := cmd.Args[0]
	currentTime := time.Now()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	if _, err := state.db.CreateUser(context.Background(), params); err != nil {
		log.Fatalf("error creating user in db: %v", err)
	}
	if err := state.cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("created and logged in as user %s\n", username)
	return nil
}

func handleUsers(state *State, cmd Command) error {
	if cmd.Name != "users" {
		return fmt.Errorf("expected users command, got %s", cmd.Name)
	}
	if len(cmd.Args) != 0 {
		return fmt.Errorf("users command expected 0 arguments, got %d arguments", len(cmd.Args))
	}
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting all users")
	}
	for _, user := range users {
		if user.Name == state.cfg.CurrentUsername {
			fmt.Println(user.Name, "(current)")
			continue
		}
		fmt.Println(user.Name)
	}
	return nil
}

func handlerReset(state *State, cmd Command) error {
	if cmd.Name != "reset" {
		return fmt.Errorf("expected reset command, got %s", cmd.Name)
	}
	if len(cmd.Args) != 0 {
		return fmt.Errorf("reset command expected 0 arguments, got %d arguments", len(cmd.Args))
	}
	if err := state.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	fmt.Println("users table reset")
	return nil
}
