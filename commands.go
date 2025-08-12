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

func handlerAgg(_ *State, cmd Command) error {
	if cmd.Name != "agg" {
		return fmt.Errorf("expected agg command, got %s", cmd.Name)
	}
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handlerAddfeed(state *State, cmd Command) error {
	if cmd.Name != "addfeed" {
		return fmt.Errorf("expected addfeed command, got %s", cmd.Name)
	}
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed command expected 2 argument (name, url), got %d", len(cmd.Args))
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	fmt.Printf("fetching feed: %s\n", url)
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	currentTime := time.Now()
	currentUserId := getCurrentUserId(state)
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      name,
		Url:       url,
		UserID:    currentUserId,
	}
	if _, err := state.db.CreateFeed(context.Background(), params); err != nil {
		return fmt.Errorf("failed to create feed entry for fedd with title: %s", feed.Channel.Title)
	}
	fmt.Printf("added feed '%s' to db as '%s'\n", feed.Channel.Title, name)
	if err = handlerFollow(state, Command{Name: "follow", Args: []string{url}}); err != nil {
		return err
	}
	return nil
}

func handlerFeeds(state *State, cmd Command) error {
	if cmd.Name != "feeds" {
		return fmt.Errorf("expected feeds command, got %s", cmd.Name)
	}
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting all users")
	}
	for _, feed := range feeds {
		feed_user, err := state.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user for feed %s with user_id=%v", feed.Name, feed.UserID)
		}
		fmt.Printf("%s / %s / %s\n", feed.Name, feed_user.Name, feed.Url)
	}
	return nil
}

func handlerFollow(state *State, cmd Command) error {
	if cmd.Name != "follow" {
		return fmt.Errorf("expected follow command, got %s", cmd.Name)
	}
	if len(cmd.Args) != 1 {
		return fmt.Errorf("follow command expected 1 argument (feed url), got %d", len(cmd.Args))
	}

	url := cmd.Args[0]

	feed, err := state.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed with url %s not found. add it with addfeed <url>", url)
	}

	currentTime := time.Now()
	user_id := getCurrentUserId(state)
	follow_params := database.CreateFeedFollowParams{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		FeedID:    feed.ID,
		UserID:    user_id,
	}

	if _, err := state.db.CreateFeedFollow(context.Background(), follow_params); err != nil {
		return err
	}
	fmt.Printf("user %s is now following %s\n", state.cfg.CurrentUsername, feed.Name)
	return nil
}

func handlerFollowing(state *State, cmd Command) error {
	user_id := getCurrentUserId(state)
	following, err := state.db.GetFeedFollowsForUser(context.Background(), user_id)
	if err != nil {
		return err
	}
	for _, feed := range following {
		fmt.Println(feed.Name)
	}
	return nil
}

func getCurrentUserId(state *State) uuid.UUID {
	user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUsername)
	if err != nil {
		log.Fatal("something went wrong trying to get the current user id. There mayb be no current user. Check with 'gator users'")
	}
	return user.ID
}

middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return nil
}
