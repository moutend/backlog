package app

import (
	"backlog/internal/cache"
	"backlog/internal/markdown"
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/moutend/go-backlog/pkg/types"

	"backlog/internal/backlog"

	"github.com/spf13/cobra"
)

var issueCommand = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"i"},
	RunE:    issueCommandRunE,
}

func issueCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var issueListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    issueListCommandRunE,
}

func issueListCommandRunE(cmd *cobra.Command, args []string) error {
	var (
		myself *types.User
		issues []*types.Issue
		ctx    context.Context
		err    error
	)

	timeout, _ := cmd.Flags().GetDuration("timeout")

	if timeout == 0 {
		goto PRINT_ISSUES
	}
	if timeout > 15*time.Minute {
		timeout = 15 * time.Minute
	}

	ctx, _ = context.WithTimeout(context.Background(), timeout)

	myself, err = backlog.GetMyself()

	if err != nil {
		// warn
	}
	if err := cache.SaveMyself(myself); err != nil {
		return err
	}

	issues, err = backlog.GetAllIssuesContext(ctx)

	if err != nil {
		// warn
	}
	if len(issues) == 0 {
		goto PRINT_ISSUES
	}
	if err := cache.Save(issues); err != nil {
		return err
	}

PRINT_ISSUES:

	myself, err = cache.LoadMyself()

	if err != nil {
		return err
	}

	issues, err = cache.LoadIssues()

	if err != nil {
		return err
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Created.Time().After(issues[j].Created.Time())
	})

	selectAssignedMe, err := cmd.Flags().GetBool("me")
	projectKey, _ := cmd.Flags().GetString("project")

	for _, issue := range issues {
		if projectKey != "" && !strings.HasPrefix(issue.IssueKey, projectKey) {
			continue
		}
		if selectAssignedMe && issue.Assignee.Id != myself.Id {
			continue
		}

		cmd.Printf(
			"- [%s] (%s) %s\n",
			issue.IssueKey,
			issue.Status.Name,
			issue.Summary,
		)
	}

	return nil
}

var issueCreateCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	RunE:    issueCreateCommandRunE,
}

func issueCreateCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	data, err := ioutil.ReadFile(args[0])

	if err != nil {
		return err
	}

	mi := &markdown.Issue{}

	if err := mi.Unmarshal(data); err != nil {
		return err
	}

	createdIssue, err := backlog.AddIssue(mi.Issue, nil)

	if err != nil {
		return err
	}

	cmd.Println("Created", createdIssue.IssueKey)

	return nil
}

var issueReadCommand = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	RunE:    issueReadCommandRunE,
}

func issueReadCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	var (
		project     *types.Project
		issue       *types.Issue
		parentIssue *types.Issue
		err         error
	)

	issue, err = backlog.GetIssue(args[0])

	if err != nil {
		goto PRINT_ISSUE
	}
	if err := cache.Save(issue); err != nil {
		return err
	}

	project, err = backlog.GetProject(fmt.Sprint(*issue.ProjectId))

	if err != nil {
		goto PRINT_ISSUE
	}
	if err := cache.Save(project); err != nil {
		return err
	}
	if issue.ParentIssueId == nil {
		goto PRINT_ISSUE
	}

	parentIssue, err = backlog.GetIssue(fmt.Sprint(*issue.ParentIssueId))

	if err != nil {
		goto PRINT_ISSUE
	}
	if err := cache.Save(parentIssue); err != nil {
		return err
	}

PRINT_ISSUE:

	issue, err = cache.LoadIssue(args[0])

	if err != nil {
		return err
	}

	project, err = cache.LoadProject(fmt.Sprint(*issue.ProjectId))

	if err != nil {
		return err
	}
	if issue.ParentIssueId != nil {
		parentIssue, err = cache.LoadIssue(fmt.Sprint(*issue.ParentIssueId))

		if err != nil {
			return err
		}
	}

	mi := &markdown.Issue{
		Project:     project,
		Issue:       issue,
		ParentIssue: parentIssue,
	}

	output, err := mi.Marshal()

	if err != nil {
		return err
	}

	cmd.Printf("%s", output)

	return nil
}

var issueUpdateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	RunE:    issueUpdateCommandRunE,
}

func issueUpdateCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	data, err := ioutil.ReadFile(args[0])

	if err != nil {
		return err
	}

	mi := &markdown.Issue{}

	if err := mi.Unmarshal(data); err != nil {
		return err
	}

	originalIssue, err := backlog.GetIssue(mi.Issue.IssueKey)

	if err != nil {
		return err
	}

	mi.Issue.Uniq(originalIssue)

	comment, _ := cmd.Flags().GetString("comment")

	updatedIssue, err := backlog.UpdateIssue(mi.Issue, nil, comment)

	if err != nil {
		return err
	}

	cmd.Println("Updated", updatedIssue.IssueKey)

	return nil
}

var issueDeleteCommand = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	RunE:    issueDeleteCommandRunE,
}

func issueDeleteCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	deletedIssue, err := backlog.DeleteIssue(args[0])

	if err != nil {
		return err
	}

	cmd.Println("Deleted", deletedIssue.IssueKey)

	return nil
}
func init() {
	RootCommand.AddCommand(issueCommand)

	issueUpdateCommand.Flags().StringP("comment", "c", "", "Set comment")

	issueListCommand.Flags().BoolP("all", "a", false, "Fetch all issues (default=false)")
	issueListCommand.Flags().BoolP("me", "m", false, "Select issues which assigned to me (default=false)")
	issueListCommand.Flags().StringP("project", "p", "", "Specify issue's project")

	issueCommand.AddCommand(issueListCommand)
	issueCommand.AddCommand(issueCreateCommand)
	issueCommand.AddCommand(issueReadCommand)
	issueCommand.AddCommand(issueUpdateCommand)
	issueCommand.AddCommand(issueDeleteCommand)
}
