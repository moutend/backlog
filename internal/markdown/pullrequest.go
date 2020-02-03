package markdown

import (
	"backlog/internal/backlog"
	"bytes"
	"fmt"

	"github.com/ericaro/frontmatter"
	"github.com/moutend/go-backlog/pkg/types"
)

type PullRequestFrontMatter struct {
	Project    string `fm:"project"`
	Summary    string `fm:"summary"`
	Repository string `fm:"repository"`
	Status     string `fm:"status"`
	Base       string `fm:"base"`
	Branch     string `fm:"branch"`
	Issue      string `fm:"issue"`
	Assignee   string `fm:"assignee"`
	Content    string `fm:"content"`
}

type PullRequest struct {
	Project     *types.Project     `json:"project"`
	PullRequest *types.PullRequest `json:"pullRequest"`
	Repository  *types.Repository  `json:"repository"`
}

func (v *PullRequest) Marshal() ([]byte, error) {
	buffer := &bytes.Buffer{}

	fmt.Fprintln(buffer, "---")

	if v.Project != nil {
		fmt.Fprintf(buffer, "project: %s\n", v.Project.ProjectKey)
	}
	if v.Repository != nil {
		fmt.Fprintf(buffer, "repository: %s\n", v.Repository.Name)
	}
	if v.PullRequest == nil {
		return nil, fmt.Errorf("markdown: PullRequest.PullRequest is required")
	}

	fmt.Fprintf(buffer, "summary: %q\n", v.PullRequest.Summary)
	fmt.Fprintf(buffer, "number: %d\n", v.PullRequest.Number)
	fmt.Fprintf(buffer, "base: %s\n", v.PullRequest.Base)
	fmt.Fprintf(buffer, "branch: %s\n", v.PullRequest.Branch)
	fmt.Fprintf(buffer, "status: %s\n", v.PullRequest.Status.Name)

	if v.PullRequest.Issue != nil {
		fmt.Fprintf(buffer, "issue: %s\n", v.PullRequest.Issue.IssueKey)
	}
	if v.PullRequest.Assignee != nil {
		fmt.Fprintf(buffer, "assignee: %s\n", v.PullRequest.Assignee.Name)
	}
	if v.PullRequest.CreatedUser != nil {
		fmt.Fprintf(buffer, "createdUser: %s\n", v.PullRequest.CreatedUser.Name)
	}
	if v.PullRequest.Created != nil {
		fmt.Fprintf(buffer, "createdAt: %s\n", v.PullRequest.Created)
	}
	if v.PullRequest.UpdatedUser != nil {
		fmt.Fprintf(buffer, "updatedUser: %s\n", v.PullRequest.UpdatedUser.Name)
	}
	if v.PullRequest.Updated != nil {
		fmt.Fprintf(buffer, "updatedAt: %s\n", v.PullRequest.Updated)
	}
	if v.PullRequest.CloseAt != nil {
		fmt.Fprintf(buffer, "closed: %s\n", v.PullRequest.CloseAt)
	}
	if v.PullRequest.MergeAt != nil {
		fmt.Fprintf(buffer, "merged: %s\n", v.PullRequest.MergeAt)
	}

	fmt.Fprintln(buffer, "---")
	fmt.Fprintln(buffer, v.PullRequest.Description)

	return buffer.Bytes(), nil
}

func (v *PullRequest) Unmarshal(data []byte) error {
	var prfm PullRequestFrontMatter

	if err := frontmatter.Unmarshal(data, &prfm); err != nil {
		return fmt.Errorf("markdown: %w", err)
	}

	project, err := backlog.GetProject(prfm.Project)

	if err != nil {
		return fmt.Errorf("markdown: %w", err)
	}

	v.Project = project

	repository, err := backlog.GetRepository(prfm.Project, prfm.Repository)

	if err != nil {
		return fmt.Errorf("markdown: %w", err)
	}

	v.Repository = repository

	v.PullRequest = &types.PullRequest{
		Summary:     prfm.Summary,
		Description: prfm.Content,
		Base:        prfm.Base,
		Branch:      prfm.Branch,
	}

	if prfm.Issue != "" {
		issue, err := backlog.GetIssue(prfm.Issue)

		if err != nil {
			return fmt.Errorf("markdown: %w", err)
		}

		v.PullRequest.Issue = issue
	}

	return nil
}
