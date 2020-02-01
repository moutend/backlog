package app

import (
	"backlog/internal/backlog"

	"github.com/spf13/cobra"
)

var userCommand = &cobra.Command{
	Use:     "user",
	Aliases: []string{"u"},
	RunE:    userCommandRunE,
}

func userCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var userListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    userListCommandRunE,
}

func userListCommandRunE(cmd *cobra.Command, args []string) error {
	users, err := backlog.GetUsers()

	if err != nil {
		return err
	}
	for _, user := range users {
		cmd.Printf("- %s (id:%d)\n", user.Name, user.Id)
	}

	return nil
}

var userReadCommand = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	RunE:    userReadCommandRunE,
}

func userReadCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	RootCommand.AddCommand(userCommand)

	userCommand.AddCommand(userListCommand)
	userCommand.AddCommand(userReadCommand)
}
