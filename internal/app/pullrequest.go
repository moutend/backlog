package app

import (
	"backlog/internal/cache"
	"backlog/internal/markdown"
	"context"
	"io/ioutil"
	"sort"
	"strconv"
	"time"

	"backlog/internal/backlog"

	"github.com/moutend/go-backlog/pkg/types"
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
	var (
		project      *types.Project
		repository   *types.Repository
		pullRequests []*types.PullRequest
		ctx          context.Context
		err          error
	)

	if len(args) < 2 {
		return nil
	}

	projectKey := args[0]
	repositoryName := args[1]

	timeout, _ := cmd.Flags().GetDuration("timeout")

	if timeout == 0 {
		goto PRINT_PULLREQUESTS
	}

	ctx, _ = context.WithTimeout(context.Background(), timeout)

	project, err = backlog.GetProject(projectKey)

	if err != nil {
		warn.Println(err)

		goto PRINT_PULLREQUESTS
	}
	if err := cache.Save(project); err != nil {
		return err
	}

	repository, err = backlog.GetRepository(projectKey, repositoryName)

	if err != nil {
		warn.Println(err)

		goto PRINT_PULLREQUESTS
	}
	if err := cache.Save(repository); err != nil {
		return err
	}

	pullRequests, err = backlog.GetAllPullRequestsContext(ctx, projectKey, repositoryName)

	if err != nil {
		warn.Println(err)
	}
	if err := cache.SavePullRequests(projectKey, repositoryName, pullRequests); err != nil {
		return err
	}

PRINT_PULLREQUESTS:

	pullRequests, err = cache.LoadPullRequests(projectKey, repositoryName)

	if err != nil {
		return err
	}

	sort.Slice(pullRequests, func(i, j int) bool {
		return pullRequests[i].Created.Time().After(pullRequests[j].Created.Time())
	})

	for _, pullRequest := range pullRequests {
		cmd.Printf(
			"%d. [%s] %s\n",
			pullRequest.Number,
			pullRequest.Status.Name,
			pullRequest.Summary)
	}

	return nil
}

var pullRequestCreateCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	RunE:    pullRequestCreateCommandRunE,
}

func pullRequestCreateCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	data, err := ioutil.ReadFile(args[0])

	if err != nil {
		return err
	}

	mpr := markdown.PullRequest{}

	mpr.Unmarshal(data)

	createdPullRequest, err := backlog.AddPullRequest(mpr.PullRequest, nil)

	if err != nil {
		return nil
	}

	cmd.Println("Created", createdPullRequest.Number)

	return nil
}

var pullRequestReadCommand = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	RunE:    pullRequestReadCommandRunE,
}

func pullRequestReadCommandRunE(cmd *cobra.Command, args []string) error {
	var (
		project     *types.Project
		repository  *types.Repository
		pullRequest *types.PullRequest
		number      int
		err         error
	)

	if len(args) < 3 {
		return nil
	}

	projectKey := args[0]
	repositoryName := args[1]
	number, err = strconv.Atoi(args[2])

	if err != nil {
		return err
	}

	project, err = backlog.GetProject(projectKey)

	if err != nil {
		goto PRINT_PULLREQUEST
	}
	if err := cache.SavePullRequest(projectKey, repositoryName, pullRequest); err != nil {
		return err
	}

	repository, err = backlog.GetRepository(project.ProjectKey, repositoryName)

	if err != nil {
		goto PRINT_PULLREQUEST
	}
	if err := cache.Save(repository); err != nil {
		return err
	}

	pullRequest, err = backlog.GetPullRequest(projectKey, repositoryName, int64(number))

	if err != nil {
		goto PRINT_PULLREQUEST
	}
	if err := cache.SavePullRequest(projectKey, repositoryName, pullRequest); err != nil {
		return err
	}

PRINT_PULLREQUEST:

	project, err = cache.LoadProject(projectKey)

	if err != nil {
		return err
	}

	repository, err = cache.LoadRepository(repositoryName)

	if err != nil {
		return err
	}

	pullRequest, err = cache.LoadPullRequest(projectKey, repositoryName, number)

	if err != nil {
		return err
	}

	mpr := markdown.PullRequest{
		Project:     project,
		Repository:  repository,
		PullRequest: pullRequest,
	}

	output, err := mpr.Marshal()

	if err != nil {
		return err
	}

	cmd.Printf("%s", output)

	return nil
}

var pullRequestUpdateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	RunE:    pullRequestUpdateCommandRunE,
}

func pullRequestUpdateCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	data, err := ioutil.ReadFile(args[0])

	if err != nil {
		return err
	}

	mpr := markdown.PullRequest{}

	mpr.Unmarshal(data)

	comment, _ := cmd.Flags().GetString("comment")

	updatedPullRequest, err := backlog.UpdatePullRequest(mpr.PullRequest, nil, comment)

	if err != nil {
		return err
	}

	cmd.Println("Updated", updatedPullRequest.Number)

	return nil
}

func init() {
	RootCommand.AddCommand(pullRequestCommand)

	pullRequestListCommand.Flags().DurationP("timeout", "t", time.Minute, "Set timeout value (default=60s)")

	pullRequestUpdateCommand.Flags().StringP("comment", "c", "", "Set comment")

	pullRequestCommand.AddCommand(pullRequestListCommand)
	pullRequestCommand.AddCommand(pullRequestReadCommand)
	pullRequestCommand.AddCommand(pullRequestUpdateCommand)
}
