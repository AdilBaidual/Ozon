package main

import (
	"Service/config"
	"Service/internal/app"
	"go.uber.org/fx"
	"log"
	"os"
	"strconv"
)

func init() {
	mode := os.Getenv("IN_MEMORY_MODE")
	if mode == "" {
		log.Fatalln("IN_MEMORY_MODE environment variable not set")
	}

	choice, err := strconv.ParseBool(mode)
	if err != nil {
		log.Fatalln("Invalid IN_MEMORY_MODE")
	}
	log.Println("IN_MEMORY_MODE ", choice)
	config.InMemory = choice
}

func main() {
	fx.New(app.NewApp()).Run()
}
