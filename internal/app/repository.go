package app

import (
	"backlog/internal/backlog"

	"github.com/moutend/go-backlog/pkg/types"

	"github.com/spf13/cobra"
)

var repositoryCommand = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"r"},
	RunE:    repositoryCommandRunE,
}

func repositoryCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var repositoryListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    repositoryListCommandRunE,
}

func repositoryListCommandRunE(cmd *cobra.Command, args []string) error {
	projects, err := backlog.GetProjects(nil)

	if err != nil {
		return err
	}

	repoMap := map[uint64][]*types.Repository{}

	for _, project := range projects {
		repos, err := backlog.GetRepositories(project.ProjectKey, nil)

		if err != nil {
			return err
		}

		repoMap[project.Id] = repos
	}
	for i, project := range projects {
		cmd.Printf("# [%s] %s\n\n", project.ProjectKey, project.Name)

		repos := repoMap[project.Id]

		if len(repos) == 0 {
			cmd.Println("No repositories.")

			goto NEXT
		}
		for _, repo := range repos {
			cmd.Printf("- %s", repo.Name)

			if yes, _ := cmd.Flags().GetBool("url"); yes {
				cmd.Printf(" (%s)", repo.HTTPURL)
			}
			if yes, _ := cmd.Flags().GetBool("ssh"); yes {
				cmd.Printf(" (%s)", repo.SSHURL)
			}

			cmd.Printf("\n")
		}

	NEXT:

		if i < len(projects)-1 {
			cmd.Println("")
		}
	}

	return nil
}

func init() {
	RootCommand.AddCommand(repositoryCommand)

	repositoryListCommand.Flags().BoolP("url", "u", false, "Show repository URL (HTTP)")
	repositoryListCommand.Flags().BoolP("ssh", "s", false, "Show repository URL (SSH)")

	repositoryCommand.AddCommand(repositoryListCommand)
}
