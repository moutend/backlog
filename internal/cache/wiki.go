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

func saveWiki(wiki *types.Wiki) error {
	if wiki == nil {
		return fmt.Errorf("cache: can't save nil wiki")
	}

	// Ensure the output directory exists.
	os.MkdirAll(cacheWikiPath, 0755)

	data, err := json.Marshal(wiki)

	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	outputPath := filepath.Join(cacheWikiPath, fmt.Sprintf("%d.json", wiki.Id))

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func saveWikis(wikis []*types.Wiki) error {
	for _, wiki := range wikis {
		if err := saveWiki(wiki); err != nil {
			return err
		}
	}

	return nil
}

func LoadWiki(wikiId uint64) (*types.Wiki, error) {
	path := filepath.Join(cacheWikiPath, fmt.Sprintf("%d.json", wikiId))
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var wiki *types.Wiki

	if err := json.Unmarshal(data, &wiki); err != nil {
		return nil, err
	}

	return wiki, nil
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

		var wiki *types.Wiki

		if err := json.Unmarshal(data, &wiki); err != nil {
			return err
		}

		wikis = append(wikis, wiki)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	return wikis, nil
}
