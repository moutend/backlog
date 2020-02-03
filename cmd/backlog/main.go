package main

import (
	"backlog/internal/app"
	"os"
)

func main() {
	if err := app.RootCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}
