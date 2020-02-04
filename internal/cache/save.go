package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/moutend/go-backlog/pkg/types"
)

func Save(v interface{}) error {
	if cachePath == "" {
		return fmt.Errorf("cache: can't save cache (probably Setup is not called)")
	}
	switch v.(type) {
	case []*types.Issue:
		return saveIssues(v.([]*types.Issue))
	case *types.Issue:
		return saveIssue(v.(*types.Issue))
	case []*types.Project:
		return saveProjects(v.([]*types.Project))
	case *types.Project:
		return saveProject(v.(*types.Project))
	}

	return fmt.Errorf("cache: type %t is not supported", v)
}

func saveIssues(issues []*types.Issue) error {
	if issues == nil {
		return fmt.Errorf("cache: can't save nil issues")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheIssuePath, 0755)

	for _, issue := range issues {
		if issue == nil {
			return fmt.Errorf("cache: can't save nil issue")
		}
		if issue.IssueKey == "" {
			return fmt.Errorf("cache: can't save issue because issue.IssueKey is empty")
		}

		data, err := json.Marshal(issue)

		if err != nil {
			return fmt.Errorf("cache: %w", err)
		}

		outputPath := filepath.Join(cacheIssuePath, issue.IssueKey+".json")

		if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
			return fmt.Errorf("cache: %w", err)
		}
	}

	return nil
}

func saveIssue(issue *types.Issue) error {
	if issue == nil {
		return fmt.Errorf("cache: issue is empty")
		return fmt.Errorf("cache: can't save nil issue")
	}
	if issue.IssueKey == "" {
		return fmt.Errorf("cache: can't save issue because issue.IssueKey is empty")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheIssuePath, 0755)

	data, err := json.Marshal(issue)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheIssuePath, issue.IssueKey+".json")

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func saveProjects(projects []*types.Project) error {
	if projects == nil {
		return fmt.Errorf("cache: can't save nil projects")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheProjectPath, 0755)

	for _, project := range projects {
		if project == nil {
			return fmt.Errorf("cache: can't save nil project")
		}
		if project.ProjectKey == "" {
			return fmt.Errorf("cache: can't save project because project.ProjectKey is empty")
		}

		data, err := json.Marshal(project)

		if err != nil {
			return fmt.Errorf("cache: %w", err)
		}

		outputPath := filepath.Join(cacheProjectPath, project.ProjectKey+".json")

		if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
			return fmt.Errorf("cache: %w", err)
		}
	}

	return nil
}

func saveProject(project *types.Project) error {
	if project == nil {
		return fmt.Errorf("cache: project is empty")
		return fmt.Errorf("cache: can't save nil project")
	}
	if project.ProjectKey == "" {
		return fmt.Errorf("cache: can't save project because project.ProjectKey is empty")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheProjectPath, 0755)

	data, err := json.Marshal(project)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheProjectPath, project.ProjectKey+".json")

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}
