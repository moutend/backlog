package app

import (
	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/backlog/internal/cache"
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
	var (
		projects     []*types.Project
		repositories []*types.Repository
		err          error
	)

	pullRequestsCountMap := map[uint64]int64{}

	timeout, _ := cmd.Flags().GetDuration("timeout")

	if timeout == 0 {
		goto PRINT_REPOSITORIES
	}

	projects, err = backlog.GetProjects(nil)

	if err != nil {
		warn.Println(err)

		goto PRINT_REPOSITORIES
	}
	if err := cache.Save(projects); err != nil {
		return err
	}
	for _, project := range projects {
		repositories, err = backlog.GetRepositories(project.ProjectKey, nil)

		if err != nil {
			warn.Println(err)

			goto PRINT_REPOSITORIES
		}
		if err := cache.Save(repositories); err != nil {
			return err
		}
		for _, repository := range repositories {
			pullRequestsCount, err := backlog.GetPullRequestsCount(project.ProjectKey, repository.Name)

			if err != nil {
				warn.Println(err)
			}

			pullRequestsCountMap[repository.Id] = pullRequestsCount
		}
	}

PRINT_REPOSITORIES:

	projects, err = cache.LoadProjects()

	if err != nil {
		return err
	}

	repositories, err = cache.LoadRepositories()

	if err != nil {
		return err
	}

	repositoriesMap := map[uint64][]*types.Repository{}

	for _, repository := range repositories {
		repositoriesMap[repository.ProjectId] = append(repositoriesMap[repository.ProjectId], repository)
	}
	for i, project := range projects {
		cmd.Printf("# [%s] %s\n\n", project.ProjectKey, project.Name)

		repositories := repositoriesMap[project.Id]

		if len(repositories) == 0 {
			cmd.Println("No repositories.")

			goto NEXT
		}
		for _, repository := range repositories {
			cmd.Printf("- %s (%d pull requests)", repository.Name, pullRequestsCountMap[repository.Id])

			if yes, _ := cmd.Flags().GetBool("url"); yes {
				cmd.Printf(" (%s)", repository.HTTPURL)
			}
			if yes, _ := cmd.Flags().GetBool("ssh"); yes {
				cmd.Printf(" (%s)", repository.SSHURL)
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
