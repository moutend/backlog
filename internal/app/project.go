package app

import (
	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/backlog/internal/cache"
	"github.com/moutend/go-backlog/pkg/types"
	"github.com/spf13/cobra"
)

var projectCommand = &cobra.Command{
	Use:     "project",
	Aliases: []string{"p"},
	RunE:    projectCommandRunE,
}

func projectCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var projectListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    projectListCommandRunE,
}

func projectListCommandRunE(cmd *cobra.Command, args []string) error {
	var (
		projects []*types.Project
		err      error
	)

	timeout, _ := cmd.Flags().GetDuration("timeout")

	if timeout == 0 {
		goto PRINT_PROJECTS
	}

	projects, err = backlog.GetProjects(nil)

	if err != nil {
		warn.Println(err)

		goto PRINT_PROJECTS
	}
	if err := cache.Save(projects); err != nil {
		return err
	}

PRINT_PROJECTS:

	projects, err = cache.LoadProjects()

	if err != nil {
		return err
	}
	for _, project := range projects {
		cmd.Printf("- [%s] %s\n", project.ProjectKey, project.Name)
	}

	return nil
}

func init() {
	RootCommand.AddCommand(projectCommand)

	projectCommand.AddCommand(projectListCommand)
}
