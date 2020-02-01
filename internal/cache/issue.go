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

func saveIssue(issue *types.Issue) error {
	if issue == nil {
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

func saveIssues(issues []*types.Issue) error {
	for _, issue := range issues {
		if err := saveIssue(issue); err != nil {
			return err
		}
	}

	return nil
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
