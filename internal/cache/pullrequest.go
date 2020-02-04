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

func SavePullRequest(projectKey, repositoryName string, pullRequest *types.PullRequest) error {
	if projectKey == "" {
		return fmt.Errorf("cache: projectKey is required")
	}
	if repositoryName == "" {
		return fmt.Errorf("cache: repositoryName is required")
	}
	if pullRequest == nil {
		return fmt.Errorf("cache: can't save nil pull request")
	}

	basePath := filepath.Join(cachePullRequestPath, projectKey, repositoryName)

	// Ensure the output directory exists.
	os.MkdirAll(basePath, 0755)

	data, err := json.Marshal(pullRequest)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(basePath, fmt.Sprintf("%d.json", pullRequest.Number))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func SavePullRequests(projectKey, repositoryName string, pullRequests []*types.PullRequest) error {
	for _, pullRequest := range pullRequests {
		if err := SavePullRequest(projectKey, repositoryName, pullRequest); err != nil {
			return err
		}
	}

	return nil
}

func LoadPullRequest(projectKey, repositoryName string, number int) (*types.PullRequest, error) {
	if projectKey == "" {
		return nil, fmt.Errorf("cache: projectKey is required")
	}
	if repositoryName == "" {
		return nil, fmt.Errorf("cache: repositoryName is required")
	}

	filePath := filepath.Join(cachePullRequestPath, projectKey, repositoryName, fmt.Sprintf("%d.json"))

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	var pullRequest *types.PullRequest

	if err := json.Unmarshal(data, &pullRequest); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return pullRequest, nil
}

func LoadPullRequests(projectKey, repositoryName string) ([]*types.PullRequest, error) {
	if projectKey == "" {
		return nil, fmt.Errorf("cache: projectKey is required")
	}
	if repositoryName == "" {
		return nil, fmt.Errorf("cache: repositoryName is required")
	}

	pullRequests := []*types.PullRequest{}
	basePath := filepath.Join(cachePullRequestPath, projectKey, repositoryName)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var pullRequest *types.PullRequest

		if err := json.Unmarshal(data, &pullRequest); err != nil {
			return err
		}

		pullRequests = append(pullRequests, pullRequest)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return pullRequests, nil
}
