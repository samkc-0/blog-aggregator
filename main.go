package main

import (
	"fmt"
	"github.com/samkc-0/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("samuel")
	fmt.Printf("DB Url: %s\n", cfg.DbUrl)
	fmt.Printf("Current User: %s\n", cfg.CurrentUsername)
}
