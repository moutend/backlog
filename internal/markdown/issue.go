package app

import (
	"backlog/internal/backlog"
	"fmt"

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
	f := IssueFrontMatter{}

	if i.Project != nil {
		f.Project = i.Project.ProjectKey
	}
	if i.ParentIssue != nil {
		f.Parent = i.ParentIssue.IssueKey
	}
	if i.Issue == nil {
		goto END_ISSUE
	}

	f.Summary = i.Issue.Summary
	f.Content = i.Issue.Description

	if i.Issue.IssueKey != "" {
		f.Issue = i.Issue.IssueKey
	}
	if i.Issue.IssueType != nil {
		f.Type = i.Issue.IssueType.Name
	}
	if i.Issue.Priority != nil {
		f.Priority = i.Issue.Priority.Name
	}
	if i.Issue.Status != nil {
		f.Status = i.Issue.Status.Name
	}
	if i.Issue.EstimatedHours != nil {
		f64 := float64(*i.Issue.EstimatedHours)
		f.Estimated = &f64
	}
	if i.Issue.ActualHours != nil {
		f64 := float64(*i.Issue.ActualHours)
		f.Actual = &f64
	}
	if i.Issue.StartDate != nil {
		f.Start = i.Issue.StartDate.Time().Format("2006-01-02")
	}
	if i.Issue.DueDate != nil {
		f.Due = i.Issue.DueDate.Time().Format("2006-01-02")
	}

END_ISSUE:

	return frontmatter.Marshal(f)
}

func (i *Issue) unmarshal(data []byte) error {
	var fo IssueFrontMatter

	if err := frontmatter.Unmarshal(data, &fo); err != nil {
		return err
	}

	project, err := backlog.GetProject(fo.Project)

	if err != nil {
		return err
	}

	i.Project = project

	i.Issue = &types.Issue{}
	i.Issue.Summary = fo.Summary
	i.Issue.Description = fo.Content

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
