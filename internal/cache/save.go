package cache

import (
	"fmt"

	"github.com/moutend/go-backlog/pkg/types"
)

func Save(v interface{}) error {
	if cachePath == "" {
		return fmt.Errorf("cache: can't save cache (probably Setup is not called)")
	}
	switch v.(type) {
	case *types.Issue:
		return saveIssue(v.(*types.Issue))
	case []*types.Issue:
		return saveIssues(v.([]*types.Issue))
	case *types.Priority:
		return savePriority(v.(*types.Priority))
	case []*types.Priority:
		return savePriorities(v.([]*types.Priority))
	case *types.Project:
		return saveProject(v.(*types.Project))
	case []*types.Project:
		return saveProjects(v.([]*types.Project))
	case *types.Repository:
		return saveRepository(v.(*types.Repository))
	case []*types.Repository:
		return saveRepositories(v.([]*types.Repository))
	case *types.User:
		return saveUser(v.(*types.User))
	case []*types.User:
		return saveUsers(v.([]*types.User))
	case *types.Wiki:
		return saveWiki(v.(*types.Wiki))
	case []*types.Wiki:
		return saveWikis(v.([]*types.Wiki))
	}

	return fmt.Errorf("cache: type %T is not supported", v)
}
