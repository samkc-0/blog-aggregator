package main

import "fmt"

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
	if err := state.cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("logged in as user %s\n", state.cfg.CurrentUsername)
	return nil
}
