package main

import (
	"log"

	"backlog/internal/app"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if err := app.RootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
