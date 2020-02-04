package cache

import (
	"fmt"
	"path/filepath"
)

const (
	cacheBaseDir = ".backlog"
)

var (
	cachePath            string
	cacheIssuePath       string
	cacheProjectPath     string
	cachePullRequestPath string
	cacheRepositoryPath  string
	cacheUserPath        string
	cacheWikiPath        string
)

func Setup(space string) error {
	if space == "" {
		return fmt.Errorf("cache: space is required")
	}

	cachePath = filepath.Join(cacheBaseDir, space)
	cacheIssuePath = filepath.Join(cacheBaseDir, space, "issue")
	cacheProjectPath = filepath.Join(cacheBaseDir, space, "project")
	cachePullRequestPath = filepath.Join(cacheBaseDir, space, "pullrequest")
	cacheRepositoryPath = filepath.Join(cacheBaseDir, space, "repository")
	cacheUserPath = filepath.Join(cacheBaseDir, space, "user")
	cacheWikiPath = filepath.Join(cacheBaseDir, space, "wiki")

	return nil
}
