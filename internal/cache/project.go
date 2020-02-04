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

func saveProjects(projects []*types.Project) error {
	if projects == nil {
		return fmt.Errorf("cache: can't save nil projects")
	}
	for _, project := range projects {
		if err := saveProject(project); err != nil {
			return err
		}
	}

	return nil
}

func saveProject(project *types.Project) error {
	if project == nil {
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

		var project *types.Project

		if err := json.Unmarshal(data, &project); err != nil {
			return err
		}

		projects = append(projects, project)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return projects, nil
}
