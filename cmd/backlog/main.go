package main

import (
	"backlog/internal/app"
	"os"
)

func main() {
	app.RootCommand.SetOutput(os.Stdout)

	if err := app.RootCommand.Execute(); err != nil {
		app.RootCommand.SetOutput(os.Stderr)

		os.Exit(-1)
	}
}
