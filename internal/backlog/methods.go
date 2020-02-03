package backlog

import (
	"context"
	"net/url"

	. "github.com/moutend/go-backlog/pkg/types"
)

func GetIssueComments(issueIdOrKey string, query url.Values) ([]*Comment, error) {
	return backlogClient.GetIssueComments(issueIdOrKey, query)
}

func GetIssueCommentsContext(ctx context.Context, issueIdOrKey string, query url.Values) ([]*Comment, error) {
	return backlogClient.GetIssueCommentsContext(ctx, issueIdOrKey, query)
}

func GetPullRequestComments(projectKeyOrId, repositoryNameOrId, number string, query url.Values) ([]*Comment, error) {
	return backlogClient.GetPullRequestComments(projectKeyOrId, repositoryNameOrId, number, query)
}

func GetPullRequestCommentsContext(ctx context.Context, projectKeyOrId, repositoryNameOrId string, number string, query url.Values) ([]*Comment, error) {
	return backlogClient.GetPullRequestCommentsContext(ctx, projectKeyOrId, repositoryNameOrId, number, query)
}

func AddIssue(issue *Issue, notifiedUsers []*User) (*Issue, error) {
	return backlogClient.AddIssue(issue, notifiedUsers)
}

func AddIssueContext(ctx context.Context, issue *Issue, notifiedUsers []*User) (*Issue, error) {
	return backlogClient.AddIssueContext(ctx, issue, notifiedUsers)
}

func GetIssue(issueKeyOrId string) (*Issue, error) {
	return backlogClient.GetIssue(issueKeyOrId)
}

func GetIssueContext(ctx context.Context, issueKeyOrId string) (*Issue, error) {
	return backlogClient.GetIssueContext(ctx, issueKeyOrId)
}

func UpdateIssue(issue *Issue, notifiedUsers []*User, comment string) (*Issue, error) {
	return backlogClient.UpdateIssue(issue, notifiedUsers, comment)
}

func UpdateIssueContext(ctx context.Context, issue *Issue, notifiedUsers []*User, comment string) (*Issue, error) {
	return backlogClient.UpdateIssueContext(ctx, issue, notifiedUsers, comment)
}

func DeleteIssue(issueKeyOrId string) (*Issue, error) {
	return backlogClient.DeleteIssue(issueKeyOrId)
}

func DeleteIssueContext(ctx context.Context, issueKeyOrId string) (*Issue, error) {
	return backlogClient.DeleteIssueContext(ctx, issueKeyOrId)
}

func GetIssues(query url.Values) ([]*Issue, error) {
	return backlogClient.GetIssues(query)
}

func GetIssuesContext(ctx context.Context, query url.Values) ([]*Issue, error) {
	return backlogClient.GetIssuesContext(ctx, query)
}

func GetIssuesCount(query url.Values) (int64, error) {
	return backlogClient.GetIssuesCount(query)
}

func GetIssuesCountContext(ctx context.Context, query url.Values) (int64, error) {
	return backlogClient.GetIssuesCountContext(ctx, query)
}

func GetIssueTypes(projectIdOrKey string) ([]*IssueType, error) {
	return backlogClient.GetIssueTypes(projectIdOrKey)
}

func GetIssueTypesContext(ctx context.Context, projectIdOrKey string) ([]*IssueType, error) {
	return backlogClient.GetIssueTypesContext(ctx, projectIdOrKey)
}

func GetNotifications(query url.Values) ([]*Notification, error) {
	return backlogClient.GetNotifications(query)
}

func GetNotificationsContext(ctx context.Context, query url.Values) ([]*Notification, error) {
	return backlogClient.GetNotificationsContext(ctx, query)
}

func GetPriorities() ([]*Priority, error) {
	return backlogClient.GetPriorities()
}

func GetPrioritiesContext(ctx context.Context) ([]*Priority, error) {
	return backlogClient.GetPrioritiesContext(ctx)
}

func GetProjectStatuses(projectIdOrKey string) ([]*ProjectStatus, error) {
	return backlogClient.GetProjectStatuses(projectIdOrKey)
}

func GetProjectStatusesContext(ctx context.Context, projectIdOrKey string) ([]*ProjectStatus, error) {
	return backlogClient.GetProjectStatusesContext(ctx, projectIdOrKey)
}

func GetProject(projectKeyOrId string) (*Project, error) {
	return backlogClient.GetProject(projectKeyOrId)
}

func GetProjectContext(ctx context.Context, projectKeyOrId string) (*Project, error) {
	return backlogClient.GetProjectContext(ctx, projectKeyOrId)
}

func GetProjects(query url.Values) ([]*Project, error) {
	return backlogClient.GetProjects(query)
}

func GetProjectsContext(ctx context.Context, query url.Values) ([]*Project, error) {
	return backlogClient.GetProjectsContext(ctx, query)
}

func AddPullRequest(projectIdOrKey, repositoryIdOrKey string, query url.Values) (*PullRequest, error) {
	return backlogClient.AddPullRequest(projectIdOrKey, repositoryIdOrKey, query)
}

func AddPullRequestContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, query url.Values) (*PullRequest, error) {
	return backlogClient.AddPullRequestContext(ctx, projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequest(projectIdOrKey, repositoryIdOrName string, number int64) (*PullRequest, error) {
	return backlogClient.GetPullRequest(projectIdOrKey, repositoryIdOrName, number)
}

func GetPullRequestContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, number int64) (*PullRequest, error) {
	return backlogClient.GetPullRequestContext(ctx, projectIdOrKey, repositoryIdOrName, number)
}

func UpdatePullRequest(projectIdOrKey, repositoryIdOrName string, number int64, query url.Values) (*PullRequest, error) {
	return backlogClient.UpdatePullRequest(projectIdOrKey, repositoryIdOrName, number, query)
}

func UpdatePullRequestContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, number int64, query url.Values) (*PullRequest, error) {
	return backlogClient.UpdatePullRequestContext(ctx, projectIdOrKey, repositoryIdOrName, number, query)
}

func GetPullRequests(projectIdOrKey, repositoryIdOrName string, query url.Values) ([]*PullRequest, error) {
	return backlogClient.GetPullRequests(projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequestsContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, query url.Values) ([]*PullRequest, error) {
	return backlogClient.GetPullRequestsContext(ctx, projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequestsCount(projectIdOrKey, repositoryIdOrName string, query url.Values) (int64, error) {
	return backlogClient.GetPullRequestsCount(projectIdOrKey, repositoryIdOrName, query)
}

func GetPullRequestsCountContext(ctx context.Context, projectIdOrKey, repositoryIdOrName string, query url.Values) (int64, error) {
	return backlogClient.GetPullRequestsCountContext(ctx, projectIdOrKey, repositoryIdOrName, query)
}

func GetRepositories(projectKeyOrId string, query url.Values) ([]*Repository, error) {
	return backlogClient.GetRepositories(projectKeyOrId, query)
}

func GetRepositoriesContext(ctx context.Context, projectKeyOrId string, query url.Values) ([]*Repository, error) {
	return backlogClient.GetRepositoriesContext(ctx, projectKeyOrId, query)
}

func GetRepository(projectKeyOrId, repositoryNameOrId string, query url.Values) (*Repository, error) {
	return backlogClient.GetRepository(projectKeyOrId, repositoryNameOrId, query)
}

func GetRepositoryContext(ctx context.Context, projectKeyOrId, repositoryNameOrId string, query url.Values) (*Repository, error) {
	return backlogClient.GetRepositoryContext(ctx, projectKeyOrId, repositoryNameOrId, query)
}

func GetSpace() (*Space, error) {
	return backlogClient.GetSpace()
}

func GetSpaceContext(ctx context.Context) (*Space, error) {
	return backlogClient.GetSpaceContext(ctx)
}

func GetUsers() ([]*User, error) {
	return backlogClient.GetUsers()
}

func GetUsersContext(ctx context.Context) ([]*User, error) {
	return backlogClient.GetUsersContext(ctx)
}

func GetMyself() (*User, error) {
	return backlogClient.GetMyself()
}

func GetMyselfContext(ctx context.Context) (*User, error) {
	return backlogClient.GetMyselfContext(ctx)
}

func GetWiki(wikiId uint64) (*Wiki, error) {
	return backlogClient.GetWiki(wikiId)
}

func GetWikiContext(ctx context.Context, wikiId uint64) (*Wiki, error) {
	return backlogClient.GetWikiContext(ctx, wikiId)
}

func GetWikis(projectIdOrKey string, query url.Values) ([]*Wiki, error) {
	return backlogClient.GetWikis(projectIdOrKey, query)
}

func GetWikisContext(ctx context.Context, projectIdOrKey string, query url.Values) ([]*Wiki, error) {
	return backlogClient.GetWikisContext(ctx, projectIdOrKey, query)
}
