package backlog

import (
	"context"
	"net/url"

	. "github.com/moutend/go-backlog/pkg/types"
)

func GetIssueComments(issueIdOrKey string, query url.Values) ([]*Comment, error) {
	return bc.GetIssueComments(issueIdOrKey, query)
}

func GetIssueCommentsContext(ctx context.Context, issueIdOrKey string, query url.Values) ([]*Comment, error) {
	return bc.GetIssueCommentsContext(ctx, issueIdOrKey, query)
}

func GetPullRequestComments(projectKeyOrId, repositoryNameOrId, number string, query url.Values) ([]*Comment, error) {
	return bc.GetPullRequestComments(projectKeyOrId, repositoryNameOrId, number, query)
}

func GetPullRequestCommentsContext(ctx context.Context, projectKeyOrId, repositoryNameOrId string, number string, query url.Values) ([]*Comment, error) {
	return bc.GetPullRequestCommentsContext(ctx, projectKeyOrId, repositoryNameOrId, number, query)
}

func AddIssue(query url.Values) (*Issue, error) {
	return bc.AddIssue(query)
}

func AddIssueContext(ctx context.Context, query url.Values) (*Issue, error) {
	return bc.AddIssueContext(ctx, query)
}

func GetIssue(issueKeyOrId string) (*Issue, error) {
	return bc.GetIssue(issueKeyOrId)
}

func GetIssueContext(ctx context.Context, issueKeyOrId string) (*Issue, error) {
	return bc.GetIssueContext(ctx, issueKeyOrId)
}

func UpdateIssue(issueKeyOrId string, query url.Values) (*Issue, error) {
	return bc.UpdateIssue(issueKeyOrId, query)
}

func UpdateIssueContext(ctx context.Context, issueKeyOrId string, query url.Values) (*Issue, error) {
	return bc.UpdateIssueContext(ctx, issueKeyOrId, query)
}

func DeleteIssue(issueKeyOrId string) (*Issue, error) {
	return bc.DeleteIssue(issueKeyOrId)
}

func DeleteIssueContext(ctx context.Context, issueKeyOrId string) (*Issue, error) {
	return bc.DeleteIssueContext(ctx, issueKeyOrId)
}

func GetIssues(query url.Values) ([]*Issue, error) {
	return bc.GetIssues(query)
}

func GetIssuesContext(ctx context.Context, query url.Values) ([]*Issue, error) {
	return bc.GetIssuesContext(ctx, query)
}

func GetIssuesCount(query url.Values) (int64, error) {
	return bc.GetIssuesCount(query)
}

func GetIssuesCountContext(ctx context.Context, query url.Values) (int64, error) {
	return bc.GetIssuesCountContext(ctx, query)
}

func GetIssueTypes(projectIdOrKey string) ([]*IssueType, error) {
	return bc.GetIssueTypes(projectIdOrKey)
}

func GetIssueTypesContext(ctx context.Context, projectIdOrKey string) ([]*IssueType, error) {
	return bc.GetIssueTypesContext(ctx, projectIdOrKey)
}

func GetNotifications(query url.Values) ([]*Notification, error) {
	return bc.GetNotifications(query)
}

func GetNotificationsContext(ctx context.Context, query url.Values) ([]*Notification, error) {
	return bc.GetNotificationsContext(ctx, query)
}

func GetPriorities() ([]*Priority, error) {
	return bc.GetPriorities()
}

func GetPrioritiesContext(ctx context.Context) ([]*Priority, error) {
	return bc.GetPrioritiesContext(ctx)
}

func GetProjectStatuses(projectIdOrKey string) ([]*ProjectStatus, error) {
	return bc.GetProjectStatuses(projectIdOrKey)
}

func GetProjectStatusesContext(ctx context.Context, projectIdOrKey string) ([]*ProjectStatus, error) {
	return bc.GetProjectStatusesContext(ctx, projectIdOrKey)
}

func GetProject(projectKeyOrId string) (*Project, error) {
	return bc.GetProject(projectKeyOrId)
}

func GetProjectContext(ctx context.Context, projectKeyOrId string) (*Project, error) {
	return bc.GetProjectContext(ctx, projectKeyOrId)
}

func GetProjects(query url.Values) ([]*Project, error) {
	return bc.GetProjects(query)
}

func GetProjectsContext(ctx context.Context, query url.Values) ([]*Project, error) {
	return bc.GetProjectsContext(ctx, query)
}

func AddPullRequest(projectIdOrKey, repositoryIdOrKey string, query url.Values) (*PullRequest, error) {
	return bc.AddPullRequest(projectIdOrKey, repositoryIdOrKey, query)
}

func AddPullRequestContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, query url.Values) (*PullRequest, error) {
	return bc.AddPullRequestContext(ctx, projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequest(projectIdOrKey, repositoryIdOrName string, number int64) (*PullRequest, error) {
	return bc.GetPullRequest(projectIdOrKey, repositoryIdOrName, number)
}

func GetPullRequestContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, number int64) (*PullRequest, error) {
	return bc.GetPullRequestContext(ctx, projectIdOrKey, repositoryIdOrName, number)
}

func UpdatePullRequest(projectIdOrKey, repositoryIdOrName string, number int64, query url.Values) (*PullRequest, error) {
	return bc.UpdatePullRequest(projectIdOrKey, repositoryIdOrName, number, query)
}

func UpdatePullRequestContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, number int64, query url.Values) (*PullRequest, error) {
	return bc.UpdatePullRequestContext(ctx, projectIdOrKey, repositoryIdOrName, number, query)
}

func GetPullRequests(projectIdOrKey, repositoryIdOrName string, query url.Values) ([]*PullRequest, error) {
	return bc.GetPullRequests(projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequestsContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, query url.Values) ([]*PullRequest, error) {
	return bc.GetPullRequestsContext(ctx, projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequestsCount(projectIdOrKey, repositoryIdOrName string, query url.Values) (int64, error) {
	return bc.GetPullRequestsCount(projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequestsCountContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, query url.Values) (int64, error) {
	return bc.GetPullRequestsCountContext(ctx, projectIdOrKey, repositoryIdOrName, query)
}

func GetRepositories(projectKeyOrId string, query url.Values) ([]*Repository, error) {
	return bc.GetRepositories(projectKeyOrId, query)
}

func GetRepositoriesContext(ctx context.Context, projectKeyOrId string, query url.Values) ([]*Repository, error) {
	return bc.GetRepositoriesContext(ctx, projectKeyOrId, query)
}

func GetRepository(projectKeyOrId, repositoryNameOrId string, query url.Values) (*Repository, error) {
	return bc.GetRepository(projectKeyOrId, repositoryNameOrId, query)
}

func GetRepositoryContext(ctx context.Context, projectKeyOrId, repositoryNameOrId string, query url.Values) (*Repository, error) {
	return bc.GetRepositoryContext(ctx, projectKeyOrId, repositoryNameOrId, query)
}

func GetSpace() (*Space, error) {
	return bc.GetSpace()
}

func GetSpaceContext(ctx context.Context) (*Space, error) {
	return bc.GetSpaceContext(ctx)
}

func GetUsers() ([]*User, error) {
	return bc.GetUsers()
}

func GetUsersContext(ctx context.Context) ([]*User, error) {
	return bc.GetUsersContext(ctx)
}

func GetMyself() (*User, error) {
	return bc.GetMyself()
}

func GetMyselfContext(ctx context.Context) (*User, error) {
	return bc.GetMyselfContext(ctx)
}

func GetWiki(wikiId uint64) (*Wiki, error) {
	return bc.GetWiki(wikiId)
}

func GetWikiContext(ctx context.Context, wikiId uint64) (*Wiki, error) {
	return bc.GetWikiContext(ctx, wikiId)
}

func GetWikis(projectIdOrKey string, query url.Values) ([]*Wiki, error) {
	return bc.GetWikis(projectIdOrKey, query)
}

func GetWikisContext(ctx context.Context, projectIdOrKey string, query url.Values) ([]*Wiki, error) {
	return bc.GetWikisContext(ctx, projectIdOrKey, query)
}
