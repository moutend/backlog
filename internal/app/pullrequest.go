package app

import (
	"backlog/internal/markdown"
	"fmt"
	"net/url"
	"strconv"

	"backlog/internal/backlog"

	"github.com/spf13/cobra"
)

var pullRequestCommand = &cobra.Command{
	Use:     "pullrequest",
	Aliases: []string{"pr"},
	RunE:    pullRequestCommandRunE,
}

func pullRequestCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var pullRequestListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    pullRequestListCommandRunE,
}

func pullRequestListCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return nil
	}

	projectKey := args[0]
	repositoryName := args[1]

	_, err := backlog.GetProject(projectKey)

	if err != nil {
		return err
	}

	query := url.Values{}

	query.Add("count", "100")

	pullRequests, err := backlog.GetPullRequests(projectKey, repositoryName, query)

	if err != nil {
		return err
	}

	for _, pullRequest := range pullRequests {
		cmd.Printf(
			"- %d [%s] %s\n",
			pullRequest.Number,
			pullRequest.Status.Name,
			pullRequest.Summary)
	}

	return nil
}

var pullRequestReadCommand = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	RunE:    pullRequestReadCommandRunE,
}

func pullRequestReadCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	i, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	wiki, err := backlog.GetWiki(uint64(i))

	if err != nil {
		return err
	}

	project, err := backlog.GetProject(fmt.Sprint(wiki.ProjectId))

	if err != nil {
		return err
	}

	mw := markdown.Wiki{
		Wiki:    wiki,
		Project: project,
	}

	output, err := mw.Marshal()

	if err != nil {
		return err
	}

	cmd.Printf("%s", output)

	return nil
}

func init() {
	RootCommand.AddCommand(pullRequestCommand)

	pullRequestCommand.AddCommand(pullRequestListCommand)
	pullRequestCommand.AddCommand(pullRequestReadCommand)
}
