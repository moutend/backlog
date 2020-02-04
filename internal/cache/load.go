package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/moutend/go-backlog/pkg/types"
)

func LoadIssues() ([]*types.Issue, error) {
	issues := []*types.Issue{}

	err := filepath.Walk(cacheIssuePath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var issue *types.Issue

		if err := json.Unmarshal(data, &issue); err != nil {
			return err
		}

		issues = append(issues, issue)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return issues, nil
}

func LoadIssue(issueIdOrKey string) (*types.Issue, error) {
	var issueId uint64

	if i64, err := strconv.ParseInt(issueIdOrKey, 10, 64); err == nil {
		issueId = uint64(i64)
	} else {
		data, err := ioutil.ReadFile(filepath.Join(cacheIssuePath, issueIdOrKey+".json"))

		if err != nil {
			return nil, fmt.Errorf("cache: %w", err)
		}

		var issue *types.Issue

		if err := json.Unmarshal(data, &issue); err != nil {
			return nil, fmt.Errorf("cache: %w", err)
		}

		return issue, nil
	}

	issues, err := LoadIssues()

	if err != nil {
		return nil, err
	}
	for _, issue := range issues {
		if issue.Id == issueId {
			return issue, nil
		}
	}

	return nil, fmt.Errorf("cache: issue not found")
}

func LoadWikis() ([]*types.Wiki, error) {
	wikis := []*types.Wiki{}

	err := filepath.Walk(cacheWikiPath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var issue *types.Wiki

		if err := json.Unmarshal(data, &issue); err != nil {
			return err
		}

		wikis = append(wikis, issue)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return wikis, nil
}

func LoadWiki(wikiKey string) (*types.Wiki, error) {
	var wiki *types.Wiki

	data, err := ioutil.ReadFile(filepath.Join(cacheWikiPath, wikiKey+".json"))

	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}
	if err := json.Unmarshal(data, &wiki); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return wiki, nil
}

func LoadProjects() ([]*types.Project, error) {
	projects := []*types.Project{}

	err := filepath.Walk(cacheProjectPath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var issue *types.Project

		if err := json.Unmarshal(data, &issue); err != nil {
			return err
		}

		projects = append(projects, issue)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return projects, nil
}

func LoadProject(projectIdOrKey string) (*types.Project, error) {
	var projectId uint64

	if i64, err := strconv.ParseInt(projectIdOrKey, 10, 64); err == nil {
		projectId = uint64(i64)
	} else {
		data, err := ioutil.ReadFile(filepath.Join(cacheProjectPath, projectIdOrKey+".json"))

		if err != nil {
			return nil, fmt.Errorf("cache: %w", err)
		}

		var project *types.Project

		if err := json.Unmarshal(data, &project); err != nil {
			return nil, fmt.Errorf("cache: %w", err)
		}

		return project, nil
	}

	projects, err := LoadProjects()

	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		if project.Id == projectId {
			return project, nil
		}
	}

	return nil, fmt.Errorf("cache: project not found")
}
