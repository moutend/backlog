package app

import (
	"backlog/internal/backlog"

	"github.com/spf13/cobra"
)

var projectCommand = &cobra.Command{
	Use:     "project",
	Aliases: []string{"p"},
	RunE:    projectCommandRunE,
}

func projectCommandRunE(cmd *cobra.Command, args []string) error {
	return projectListCommandRunE(cmd, args)
}

var projectListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    projectListCommandRunE,
}

func projectListCommandRunE(cmd *cobra.Command, args []string) error {
	projects, err := backlog.GetProjects(nil)

	if err != nil {
		return err
	}

	for _, project := range projects {
		cmd.Printf("- %s (id:%d)\n", project.Name, project.Id)
	}

	return nil
}

func init() {
	RootCommand.AddCommand(projectCommand)

	projectCommand.AddCommand(projectListCommand)
}
