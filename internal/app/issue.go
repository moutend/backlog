package app

import (
	"backlog/internal/cache"
	"backlog/internal/markdown"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
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
		myself          *types.User
		projects        []*types.Project
		projectStatuses []*types.ProjectStatus
		issues          []*types.Issue

		projectStatusesMap map[string][]*types.ProjectStatus
		query              url.Values

		ctx context.Context
		err error
	)

	maxIssues, _ := cmd.Flags().GetInt("count")

	if maxIssues == 0 {
		return nil
	}
	if yes, _ := cmd.Flags().GetBool("all"); yes {
		maxIssues = 0
	}

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
		warn.Println(err)

		goto PRINT_ISSUES
	}
	if err := cache.SaveMyself(myself); err != nil {
		return err
	}

	projects, err = backlog.GetProjects(nil)

	if err != nil {
		warn.Println(err)

		goto PRINT_ISSUES
	}
	if err := cache.Save(projects); err != nil {
		return err
	}

	projectStatusesMap = map[string][]*types.ProjectStatus{}

	for _, project := range projects {
		projectStatuses, err = backlog.GetProjectStatuses(project.ProjectKey)

		if err != nil {
			warn.Println(err)

			goto PRINT_ISSUES
		}
		if err := cache.SaveProjectStatuses(project.ProjectKey, projectStatuses); err != nil {
			return err
		}

		projectStatusesMap[fmt.Sprint(project.Id)] = projectStatuses
	}

	query = url.Values{}

	{
		desc, _ := cmd.Flags().GetBool("desc")
		asc, _ := cmd.Flags().GetBool("desc")

		if !desc && asc {
			query.Add("order", "asc")
		} else {
			query.Add("order", "desc")
		}
	}
	if yes, _ := cmd.Flags().GetBool("myself"); yes {
		query.Add("assigneeId", fmt.Sprint(myself.Id))
	}
	if priority, _ := cmd.Flags().GetString("priority"); priority != "" {
		// query.Add("priorityId", fmt.Sprint(priority.Id))
	}
	if projectKey, _ := cmd.Flags().GetString("project"); projectKey != "" {
		for _, project := range projects {
			if projectKey == project.ProjectKey {
				query.Add("projectId", fmt.Sprint(project.Id))

				break
			}
		}
	}
	if status, _ := cmd.Flags().GetString("status"); status != "" && query.Get("projectId") != "" {
		for _, projectStatus := range projectStatusesMap[query.Get("projectId")] {
			if projectStatus.Name == status {
				query.Add("statusId", fmt.Sprint(projectStatus.Id))

				break
			}
		}
	}
	if sortedBy, _ := cmd.Flags().GetString("sort"); sortedBy != "" {
		query.Add("sort", sortedBy)
	}

	issues, err = backlog.GetAllIssuesContext(ctx, maxIssues, query)

	if err != nil {
		warn.Println(err)
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
	projects, err = cache.LoadProjects()

	if err != nil {
		return err
	}

	issues, err = cache.LoadIssues()

	if err != nil {
		return err
	}

	sortedBy, _ := cmd.Flags().GetString("order")

	switch strings.ToLower(sortedBy) {
	case "created":
		sort.Slice(issues, func(i, j int) bool {
			return issues[i].Created.Time().After(issues[j].Created.Time())
		})
	case "updated":
		sort.Slice(issues, func(i, j int) bool {
			return issues[i].Updated.Time().After(issues[j].Updated.Time())
		})
	}

	selectAssignedMe, err := cmd.Flags().GetBool("me")
	projectKey, _ := cmd.Flags().GetString("project")
	issueStatus, _ := cmd.Flags().GetString("status")
	issueCount := 0

	for _, issue := range issues {
		if issueCount == maxIssues {
			break
		}
		if issueStatus != "" && issue.Status.Name != issueStatus {
			continue
		}
		if projectKey != "" && !strings.HasPrefix(issue.IssueKey, projectKey) {
			continue
		}
		if selectAssignedMe {
			if issue.Assignee == nil {
				continue
			}
			if issue.Assignee.Id != myself.Id {
				continue
			}
		}

		cmd.Printf(
			"- [%s] (%s) %s\n",
			issue.IssueKey,
			issue.Status.Name,
			issue.Summary,
		)

		issueCount += 1
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
	issueListCommand.Flags().IntP("count", "c", 20, "Print issues (default=20)")
	issueListCommand.Flags().BoolP("desc", "", true, "Print issues descending order (default=true)")
	issueListCommand.Flags().BoolP("asc", "", false, "Print issues ascending order (default=false)")
	issueListCommand.Flags().BoolP("myself", "", false, "Select issues which assignee is myself")
	issueListCommand.Flags().StringP("sort", "", "created", "Specify sorting order (default='created')")
	issueListCommand.Flags().StringP("project", "", "", "Specify issue's project key (default='')")
	issueListCommand.Flags().StringP("status", "", "", "Specify issue status (default='')")
	issueListCommand.Flags().StringP("priority", "", "", "Specify issue priority (default='')")

	issueCommand.AddCommand(issueListCommand)
	issueCommand.AddCommand(issueCreateCommand)
	issueCommand.AddCommand(issueReadCommand)
	issueCommand.AddCommand(issueUpdateCommand)
	issueCommand.AddCommand(issueDeleteCommand)
}
