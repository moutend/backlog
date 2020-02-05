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

func saveRepository(repository *types.Repository) error {
	if repository == nil {
		return fmt.Errorf("cache: can't save nil repository")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheRepositoryPath, 0755)

	data, err := json.Marshal(repository)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheRepositoryPath, fmt.Sprintf("%s.json", repository.Name))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func saveRepositories(repositories []*types.Repository) error {
	for _, repository := range repositories {
		if err := saveRepository(repository); err != nil {
			return err
		}
	}

	return nil
}

func LoadRepository(repositoryName string) (*types.Repository, error) {
	data, err := ioutil.ReadFile(filepath.Join(cacheRepositoryPath, repositoryName+".json"))

	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	var repository *types.Repository

	if err := json.Unmarshal(data, &repository); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return repository, nil
}

func LoadRepositories() ([]*types.Repository, error) {
	repositories := []*types.Repository{}

	err := filepath.Walk(cacheRepositoryPath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var repository *types.Repository

		if err := json.Unmarshal(data, &repository); err != nil {
			return err
		}

		repositories = append(repositories, repository)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return repositories, nil
}
