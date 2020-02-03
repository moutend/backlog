package app

import (
	"backlog/internal/markdown"
	"fmt"
	"io/ioutil"

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
	issues, err := backlog.GetIssues(nil)

	if err != nil {
		return err
		// warn.Printf("Failed to get issues")
	}

	createdOrUpdatedUser := ""
	createdOrUpdatedDate := ""

	for _, issue := range issues {
		cmd.Printf(
			"- [%s] (%s) %s (updated at %s by %s)\n",
			issue.IssueKey,
			issue.Status.Name,
			issue.Summary,
			createdOrUpdatedDate,
			createdOrUpdatedUser,
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

	filePath := args[0]

	data, err := ioutil.ReadFile(filePath)

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

	mi := &markdown.Issue{}

	issue, err := backlog.GetIssue(args[0])

	if err != nil {
		return err
	}
	if issue.ProjectId == nil {
		return fmt.Errorf("issue is broken")
	}

	project, err := backlog.GetProject(fmt.Sprint(*issue.ProjectId))

	if err != nil {
		return err
	}

	mi.Issue = issue
	mi.Project = project

	if issue.ParentIssueId != nil {
		parent, err := backlog.GetIssue(fmt.Sprint(issue.ParentIssueId))

		if err != nil {
			return err
		}

		mi.ParentIssue = parent
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

	issueCommand.AddCommand(issueListCommand)
	issueCommand.AddCommand(issueCreateCommand)
	issueCommand.AddCommand(issueReadCommand)
	issueCommand.AddCommand(issueUpdateCommand)
	issueCommand.AddCommand(issueDeleteCommand)
}
