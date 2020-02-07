package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/moutend/go-backlog/pkg/types"
)

func SaveIssueComment(issueKey string, comment *types.Comment) error {
	if comment == nil {
		return fmt.Errorf("cache: can't save nil comment")
	}
	if issueKey == "" {
		return fmt.Errorf("cache: can't save comment because issueKey is empty")
	}

	// Ensure the output directory exists.
	basePath := filepath.Join(cacheIssueCommentPath, issueKey)
	os.MkdirAll(basePath, 0755)

	data, err := json.Marshal(comment)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(basePath, fmt.Sprintf("%d.json", comment.Id))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func SaveIssueComments(issueKey string, comments []*types.Comment) error {
	for _, comment := range comments {
		if err := SaveIssueComment(issueKey, comment); err != nil {
			return err
		}
	}

	return nil
}

func LoadIssueComments(issueKey string) ([]*types.Comment, error) {
	comments := []*types.Comment{}

	err := filepath.Walk(filepath.Join(cacheIssueCommentPath, issueKey), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var comment *types.Comment

		if err := json.Unmarshal(data, &comment); err != nil {
			return err
		}

		comments = append(comments, comment)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return comments, nil
}

func SavePullRequestComment(projectKey, repositoryName string, number int64, comment *types.Comment) error {
	if comment == nil {
		return fmt.Errorf("cache: can't save nil comment")
	}
	if projectKey == "" {
		return fmt.Errorf("cache: can't save comment because projectKey is empty")
	}
	if repositoryName == "" {
		return fmt.Errorf("cache: can't save comment because repositoryName is empty")
	}
	if number < 0 {
		return fmt.Errorf("cache: can't save comment because number is invalid")
	}

	// Ensure the output directory exists.
	basePath := filepath.Join(cachePullRequestCommentPath, projectKey, repositoryName, fmt.Sprint(number))
	os.MkdirAll(basePath, 0755)

	data, err := json.Marshal(comment)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(basePath, fmt.Sprintf("%d.json", comment.Id))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func SavePullRequestComments(projectKey, repositoryName string, number int64, comments []*types.Comment) error {
	for _, comment := range comments {
		if err := SavePullRequestComment(projectKey, repositoryName, number, comment); err != nil {
			return err
		}
	}

	return nil
}

func LoadPullRequestComments(projectKey, repositoryName string, number int64) ([]*types.Comment, error) {
	comments := []*types.Comment{}

	err := filepath.Walk(filepath.Join(cachePullRequestCommentPath, projectKey, repositoryName, fmt.Sprint(number)), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var comment *types.Comment

		if err := json.Unmarshal(data, &comment); err != nil {
			return err
		}

		comments = append(comments, comment)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return comments, nil
}
