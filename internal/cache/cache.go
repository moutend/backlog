package cache

import (
	"fmt"
	"path/filepath"
)

const (
	cacheBaseDir = ".backlog"
)

var (
	cachePath              string
	cacheIssuePath         string
	cachePriorityPath      string
	cacheProjectPath       string
	cacheProjectStatusPath string
	cachePullRequestPath   string
	cacheRepositoryPath    string
	cacheUserPath          string
	cacheWikiPath          string
)

func Setup(space string) error {
	if space == "" {
		return fmt.Errorf("cache: space is required")
	}

	cachePath = filepath.Join(cacheBaseDir, space)
	cacheIssuePath = filepath.Join(cacheBaseDir, space, "issue")
	cachePriorityPath = filepath.Join(cacheBaseDir, space, "priority")
	cacheProjectPath = filepath.Join(cacheBaseDir, space, "project")
	cacheProjectStatusPath = filepath.Join(cacheBaseDir, space, "projectstatus")
	cachePullRequestPath = filepath.Join(cacheBaseDir, space, "pullrequest")
	cacheRepositoryPath = filepath.Join(cacheBaseDir, space, "repository")
	cacheUserPath = filepath.Join(cacheBaseDir, space, "user")
	cacheWikiPath = filepath.Join(cacheBaseDir, space, "wiki")

	return nil
}
