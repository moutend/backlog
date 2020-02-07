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

func savePriority(priority *types.Priority) error {
	if priority == nil {
		return fmt.Errorf("cache: can't save nil priority")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cachePriorityPath, 0755)

	data, err := json.Marshal(priority)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cachePriorityPath, fmt.Sprintf("%d.json", priority.Id))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func savePriorities(priorities []*types.Priority) error {
	for _, priority := range priorities {
		if err := savePriority(priority); err != nil {
			return err
		}
	}

	return nil
}

func LoadPriorities() ([]*types.Priority, error) {
	priorities := []*types.Priority{}

	err := filepath.Walk(cachePriorityPath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		var priority *types.Priority

		if err := json.Unmarshal(data, &priority); err != nil {
			return err
		}

		priorities = append(priorities, priority)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return priorities, nil
}
