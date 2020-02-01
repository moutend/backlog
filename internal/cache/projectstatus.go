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

func SaveProjectStatus(projectKey string, projectStatus *types.ProjectStatus) error {
	if projectKey == "" {
		return fmt.Errorf("cache: can't save project status because projectKey is empty")
	}
	if projectStatus == nil {
		return fmt.Errorf("cache: can't save nil project status")
	}

	// Ensure the output directory exists.
	os.MkdirAll(filepath.Join(cacheProjectStatusPath, projectKey), 0755)

	data, err := json.Marshal(projectStatus)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheProjectStatusPath, projectKey, fmt.Sprintf("%s.json", projectStatus.Id))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func SaveProjectStatuses(projectKey string, projectStatuses []*types.ProjectStatus) error {
	for _, projectStatus := range projectStatuses {
		if err := SaveProjectStatus(projectKey, projectStatus); err != nil {
			return err
		}
	}

	return nil
}

func LoadProjectStatuses(projectKey string) ([]*types.ProjectStatus, error) {
	if projectKey == "" {
		return nil, fmt.Errorf("cache: can't save project status because projectKey is empty")
	}
	projectStatuses := []*types.ProjectStatus{}

	err := filepath.Walk(filepath.Join(cacheProjectStatusPath, projectKey), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var projectStatus *types.ProjectStatus

		if err := json.Unmarshal(data, &projectStatus); err != nil {
			return err
		}

		projectStatuses = append(projectStatuses, projectStatus)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return projectStatuses, nil
}
