package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/moutend/go-backlog/pkg/types"
)

func SaveMyself(user *types.User) error {
	if user == nil {
		return fmt.Errorf("cache: can't save nil user")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheUserPath, 0755)

	data, err := json.Marshal(user)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheUserPath, "myself.json")

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func LoadMyself() (*types.User, error) {
	data, err := ioutil.ReadFile(filepath.Join(cacheUserPath, "myself.json"))

	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	var myself *types.User

	if err := json.Unmarshal(data, &myself); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return myself, nil
}

func saveUser(user *types.User) error {
	if user == nil {
		return fmt.Errorf("cache: can't save nil user")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheUserPath, 0755)

	data, err := json.Marshal(user)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheUserPath, fmt.Sprintf("%s.json", user.Id))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func saveUsers(users []*types.User) error {
	for _, user := range users {
		if err := saveUser(user); err != nil {
			return err
		}
	}

	return nil
}
