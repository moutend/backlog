package markdown

import (
	"bytes"
	"fmt"

	"github.com/moutend/backlog/internal/backlog"

	"github.com/ericaro/frontmatter"
	"github.com/moutend/go-backlog/pkg/types"
)

type IssueFrontMatter struct {
	Summary   string   `fm:"summary"`
	Project   string   `fm:"project"`
	Issue     string   `fm:"issue"`
	Parent    string   `fm:"parent"`
	Type      string   `fm:"type"`
	Priority  string   `fm:"priority"`
	Status    string   `fm:"status"`
	Start     string   `fm:"start"`
	Due       string   `fm:"due"`
	Estimated *float64 `fm:"estimated"`
	Actual    *float64 `fm:"actual"`
	Assignee  string   `fm:"assignee"`
	Created   string   `fm:"created"`
	Updated   string   `fm:"updated"`
	Content   string   `fm:"content"`
}

type Issue struct {
	Project     *types.Project `json:"project"`
	ParentIssue *types.Issue   `json:"parentIssue"`
	Issue       *types.Issue   `json:"issue"`
}

func (i *Issue) Marshal() ([]byte, error) {
	buffer := &bytes.Buffer{}

	fmt.Fprintln(buffer, "---")

	if i.Project != nil {
		fmt.Fprintf(buffer, "project: %s\n", i.Project.ProjectKey)
	}
	if i.ParentIssue != nil {
		fmt.Fprintf(buffer, "parent: %s\n", i.ParentIssue.IssueKey)
	}
	if i.Issue == nil {
		goto END_ISSUE
	}

	fmt.Fprintf(buffer, "summary: %q\n", i.Issue.Summary)

	if i.Issue.IssueKey != "" {
		fmt.Fprintf(buffer, "issue: %s\n", i.Issue.IssueKey)
	}
	if i.Issue.IssueType != nil {
		fmt.Fprintf(buffer, "type: %s\n", i.Issue.IssueType.Name)
	}
	if i.Issue.Priority != nil {
		fmt.Fprintf(buffer, "priority: %s\n", i.Issue.Priority.Name)
	}
	if i.Issue.Status != nil {
		fmt.Fprintf(buffer, "status: %s\n", i.Issue.Status.Name)
	}
	if i.Issue.EstimatedHours != nil {
		fmt.Fprintf(buffer, "estimated: %v\n", *i.Issue.EstimatedHours)
	}
	if i.Issue.ActualHours != nil {
		fmt.Fprintf(buffer, "actual: %v\n", *i.Issue.ActualHours)
	}
	if i.Issue.StartDate != nil {
		fmt.Fprintf(buffer, "start: %s\n", i.Issue.StartDate.Time().Format("2006-01-02"))
	}
	if i.Issue.DueDate != nil {
		fmt.Fprintf(buffer, "due: %s\n", i.Issue.DueDate.Time().Format("2006-01-02"))
	}
	if i.Issue.Assignee != nil {
		fmt.Fprintf(buffer, "assignee: %s\n", i.Issue.Assignee.Name)
	}
	if i.Issue.Created != nil {
		fmt.Fprintf(buffer, "created: %s by %s\n", i.Issue.Created.Time().Format("2006-01-02"), i.Issue.CreatedUser.Name)
	}
	if i.Issue.Updated != nil {
		fmt.Fprintf(buffer, "updated: %s by %s\n", i.Issue.Updated.Time().Format("2006-01-02"), i.Issue.UpdatedUser.Name)
	}
	if i.Issue.IssueKey != "" {
		fmt.Fprintf(buffer, "url: https://%s/view/%s\n", backlog.SpaceName(), i.Issue.IssueKey)
	}

	fmt.Fprintln(buffer, "---")
	fmt.Fprint(buffer, i.Issue.Description)

END_ISSUE:

	return buffer.Bytes(), nil
}

func (i *Issue) unmarshal(data []byte) error {
	var fo IssueFrontMatter

	if err := frontmatter.Unmarshal(data, &fo); err != nil {
		return err
	}
	if fo.Project == "" {
		return fmt.Errorf("markdown: project is required")
	}

	project, err := backlog.GetProject(fo.Project)

	if err != nil {
		return err
	}

	i.Project = project

	i.Issue = &types.Issue{}

	i.Issue.Summary = fo.Summary
	i.Issue.Description = fo.Content
	i.Issue.IssueKey = fo.Issue
	i.Issue.ProjectId = &project.Id

	if fo.Estimated != nil {
		i.Issue.EstimatedHours = types.NewHours(*fo.Estimated)
	}
	if fo.Actual != nil {
		i.Issue.ActualHours = types.NewHours(*fo.Actual)
	}
	if fo.Start != "" {
		i.Issue.StartDate = types.NewDate(fo.Start)
	}
	if fo.Due != "" {
		i.Issue.DueDate = types.NewDate(fo.Due)
	}
	if fo.Parent != "" {
		parent, err := backlog.GetIssue(fo.Parent)

		if err != nil {
			return err
		}

		i.ParentIssue = parent

		i.Issue.ParentIssueId = &parent.Id
	}

	issueTypes, err := backlog.GetIssueTypes(project.ProjectKey)

	if err != nil {
		return err
	}
	for _, issueType := range issueTypes {
		if fo.Type == issueType.Name {
			i.Issue.IssueType = issueType

			break
		}
	}
	if i.Issue.IssueType == nil {
		return fmt.Errorf("markdown: invalid issue type")
	}

	priorities, err := backlog.GetPriorities()

	if err != nil {
		return err
	}
	for _, priority := range priorities {
		if fo.Priority == priority.Name {
			i.Issue.Priority = priority

			break
		}
	}
	if i.Issue.Priority == nil {
		return fmt.Errorf("markdown: invalid priority")
	}

	projectStatuses, err := backlog.GetProjectStatuses(project.ProjectKey)

	if err != nil {
		return err
	}
	for _, projectStatus := range projectStatuses {
		if fo.Status == projectStatus.Name {
			i.Issue.Status = projectStatus

			break
		}
	}
	if i.Issue.Status == nil {
		return fmt.Errorf("markdown: invalid status")
	}

	return nil
}

func (i *Issue) Unmarshal(data []byte) error {
	return i.unmarshal(data)
}
