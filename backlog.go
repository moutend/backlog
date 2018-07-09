package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ericaro/frontmatter"
	"github.com/moutend/go-backlog"
)

var (
	version   = "dev"
	revision  = "latest"
	spaceName string
	client    *backlog.Client
)

func parseMarkdown(filename string) (url.Values, error) {
	var wg sync.WaitGroup

	priorityNameToId := make(map[string]int)
	projectNameToId := make(map[string]int)
	issueTypeNameToId := make(map[string]map[string]int)
	statusNameToId := make(map[string]int)
	myselfId := ""

	wg.Add(1)
	go func() {
		statuses, _ := client.GetStatuses()
		for _, status := range statuses {
			statusNameToId[status.Name] = status.Id
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		priorities, _ := client.GetPriorities()
		for _, priority := range priorities {
			priorityNameToId[priority.Name] = priority.Id
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		projects, _ := client.GetProjects(nil)
		for _, project := range projects {
			projectNameToId[project.Name] = project.Id
			issueTypes, _ := client.GetIssueTypes(project.Id)
			issueTypeNameToId[project.Name] = make(map[string]int)
			for _, issueType := range issueTypes {
				issueTypeNameToId[project.Name][issueType.Name] = issueType.Id
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		myself, _ := client.GetMyself()
		myselfId = strconv.Itoa(myself.Id)
		wg.Done()
	}()

	wg.Wait()

	type frontmatterOption struct {
		Project       string `fm:"project"`
		ProjectId     string `fm:"projectid"`
		IssueType     string `fm:"issuetype"`
		IssueTypeId   string `fm:"issuetypeid"`
		Priority      string `fm:"priority"`
		PriorityId    string `fm:"priorityid"`
		Status        string `fm:"status"`
		StatusId      string `fm:"statusid"`
		ParentIssueId string `fm:"parentissueid"`
		Estimated     string `fm:"estimated"`
		Actual        string `fm:"actual"`
		Due           string `fm:"due"`
		Summary       string `fm:"summary"`
		Description   string `fm:"content"`
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fo := &frontmatterOption{}
	err = frontmatter.Unmarshal(file, fo)
	if err != nil {
		return nil, err
	}

	values := url.Values{}

	if fo.ProjectId != "" {
		values.Add("projectId", fo.ProjectId)
	} else if fo.Project != "" {
		values.Add("projectId", strconv.Itoa(int(projectNameToId[fo.Project])))
	} else {
		return nil, fmt.Errorf("specify project or projectid")
	}
	if fo.IssueTypeId != "" {
		values.Add("issueTypeId", fo.IssueTypeId)
	} else if fo.IssueType != "" {
		values.Add("issueTypeId", strconv.Itoa(issueTypeNameToId[fo.Project][fo.IssueType]))
	} else {
		return nil, fmt.Errorf("specify type or typeid")
	}
	if fo.PriorityId != "" {
		values.Add("priorityId", fo.PriorityId)
	} else if fo.Priority != "" {
		values.Add("priorityId", strconv.Itoa(priorityNameToId[fo.Priority]))
	} else {
		return nil, fmt.Errorf("specify priority or priorityid ")
	}
	if fo.StatusId != "" {
		values.Add("statusId", fo.StatusId)
	} else if fo.Status != "" {
		values.Add("statusId", strconv.Itoa(statusNameToId[fo.Status]))
	} else {
		return nil, fmt.Errorf("specify status or statusid")
	}
	if fo.ParentIssueId != "" {
		values.Add("parentIssueId", fo.ParentIssueId)
	}
	if fo.Estimated != "" {
		values.Add("estimatedHours", fo.Estimated)
	}
	if fo.Actual != "" {
		values.Add("actualHours", fo.Actual)
	}
	if fo.Due != "" {
		values.Add("dueDate", fo.Due)
	}

	values.Add("assigneeId", myselfId)
	values.Add("startDate", time.Now().Format("2006-01-02"))
	values.Add("summary", fo.Summary)
	values.Add("description", fo.Description)
	return values, nil
}

func parseMarkdownForPR(filename string) (url.Values, error) {
	var wg sync.WaitGroup

	projectNameToId := make(map[string]int)
	myselfId := ""

	wg.Add(1)
	go func() {
		projects, _ := client.GetProjects(nil)
		for _, project := range projects {
			projectNameToId[project.Name] = project.Id
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		myself, _ := client.GetMyself()
		myselfId = strconv.Itoa(myself.Id)
		wg.Done()
	}()

	wg.Wait()

	type frontmatterOption struct {
		Summary     string `fo:"summary"`
		IssueKey    string `fo:"issuekey"`
		Repository  string `fo:"repository"`
		Base        string `fo:"base"`
		Branch      string `fo:"branch"`
		Project     string `fm:"project"`
		ProjectId   string `fm:"projectid"`
		Description string `fm:"content"`
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fo := &frontmatterOption{}
	err = frontmatter.Unmarshal(file, fo)
	if err != nil {
		return nil, err
	}

	values := url.Values{}

	if fo.ProjectId != "" {
		values.Add("projectId", fo.ProjectId)
	} else if fo.Project != "" {
		values.Add("projectId", strconv.Itoa(int(projectNameToId[fo.Project])))
	} else {
		return nil, fmt.Errorf("specify project or projectid")
	}
	values.Add("assigneeId", myselfId)
	values.Add("summary", fo.Summary)
	values.Add("branch", fo.Branch)
	values.Add("base", fo.Base)
	values.Add("repositoryId", fo.Repository)
	values.Add("description", fo.Description)

	issue, err := client.GetIssue(fo.IssueKey)
	if fo.IssueKey != "" && err != nil {
		return nil, err
	}
	if fo.IssueKey != "" {
		values.Add("issueId", fmt.Sprintf("%v", issue.Id))
	}

	return values, nil
}

func main() {
	var err error
	if err = run(os.Args); err != nil {
		log.Fatal(err)
	}
	return
}

func run(args []string) error {
	if len(args) < 2 {
		return HelpCommand(args)
	}
	switch args[1] {
	case "v":
		return VersionCommand(args)
	case "version":
		return VersionCommand(args)
	case "h":
		return HelpCommand(args)
	case "help":
		return HelpCommand(args)
	}

	var err error
	var debugFlag bool

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[1:])
	args = f.Args()
	command := args[0]
	args = args[1:]
	spaceName = os.Getenv("BACKLOG_SPACE")
	if client, err = backlog.New(spaceName, os.Getenv("BACKLOG_TOKEN")); err != nil {
		return err
	}
	if debugFlag {
		client.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	switch command {
	case "l":
		return ListCommand(args)
	case "list":
		return ListCommand(args)
	case "lpr":
		return ListPullRequestsCommand(args)
	case "gpr":
		return GetPullRequestCommand(args)
	case "apr":
		return AddPullRequestCommand(args)
	case "upr":
		return UpdatePullRequestCommand(args)
	case "p":
		return CreateIssueCommand(args)
	case "post":
		return CreateIssueCommand(args)
	case "d":
		return DeleteIssueCommand(args)
	case "delete":
		return DeleteIssueCommand(args)
	case "g":
		return GetIssueCommand(args)
	case "get":
		return GetIssueCommand(args)
	case "u":
		return UpdateIssueCommand(args)
	case "update":
		return UpdateIssueCommand(args)
	case "c":
		return GetCommentsCommand(args)
	case "comments":
		return GetCommentsCommand(args)
	default:
		return fmt.Errorf("%s is not a subcommand", command)
	}
}

func ListPullRequestsCommand(args []string) error {
	f := flag.NewFlagSet("list", flag.ExitOnError)
	f.Parse(args)
	args = f.Args()

	projects, err := client.GetProjects(nil)
	if err != nil {
		return nil
	}
	for _, project := range projects {
		fmt.Printf("- %v (id:%v)\n", project.Name, project.Id)
		projectId := fmt.Sprintf("%v", project.Id)
		repos, err := client.GetRepositories(projectId, nil)
		if err != nil {
			return err
		}
		for _, repo := range repos {
			fmt.Println(repo.Name)
			query := url.Values{}
			query.Add("count", "100")
			repositoryId := fmt.Sprintf("%v", repo.Id)
			count, err := client.GetPullRequestsCount(projectId, repositoryId, nil)
			if err != nil {
				return err
			}
			fmt.Println(count)
			pullRequests, err := client.GetPullRequests(projectId, repositoryId, query)
			if err != nil {
				return err
			}
			for _, pullRequest := range pullRequests {
				fmt.Printf("[%v #%v] %v (%v -> %v)\n", pullRequest.Status.Name, pullRequest.Number, pullRequest.Summary, pullRequest.Branch, pullRequest.Base)
			}
		}
	}

	return nil
}

func ListCommand(args []string) error {
	var isAssignedToMe bool
	var statusFlag string
	var myself *backlog.User

	f := flag.NewFlagSet("list", flag.ExitOnError)
	f.BoolVar(&isAssignedToMe, "assigned-to-me", false, "get tasks which assigned to me")
	f.BoolVar(&isAssignedToMe, "m", false, "get tasks which assigned to me")
	f.StringVar(&statusFlag, "s", "", "specify status")
	f.Parse(args)
	args = f.Args()

	myself, err := client.GetMyself()
	if err != nil {
		return err
	}
	statuses, err := client.GetStatuses()
	if err != nil {
		return err
	}
	statusMap := make(map[string]int)
	for _, status := range statuses {
		statusMap[status.Name] = status.Id
	}
	projects, err := client.GetProjects(nil)
	if err != nil {
		return nil
	}
	for _, project := range projects {
		fmt.Printf("- %v (id:%v)\n", project.Name, project.Id)

		query := url.Values{}
		query.Add("projectId[]", fmt.Sprintf("%v", project.Id))
		query.Add("count", "100")
		if isAssignedToMe {
			query.Add("assigneeId[]", fmt.Sprintf("%v", myself.Id))
		}
		if statusID, ok := statusMap[statusFlag]; ok {
			query.Add("statusId[]", fmt.Sprintf("%v", statusID))
		}
		query.Add("sort", "updated")
		issues, err := client.GetIssues(query)
		if err != nil {
			return err
		}
		for _, issue := range issues {
			fmt.Printf("  - [%v %v] %v by @%v (id:%v) %v\n", issue.Status.Name, issue.IssueType.Name, issue.Summary, issue.CreatedUser.Name, issue.IssueKey, issue.StartDate)
		}
	}

	return nil
}

func CreateIssueCommand(args []string) error {
	if len(args) < 1 {
		return nil
	}

	values, err := parseMarkdown(args[0])
	if err != nil {
		return err
	}
	delete(values, "statusId")
	issue, err := client.CreateIssue(values)
	if err != nil {
		return err
	}

	fmt.Println(issue.Id)

	return nil
}

func AddPullRequestCommand(args []string) error {
	if len(args) < 1 {
		return nil
	}

	values, err := parseMarkdownForPR(args[0])
	if err != nil {
		return err
	}
	projectId := values.Get("projectId")
	repositoryId := values.Get("repositoryId")
	values.Del("projectId")
	values.Del("repositoryId")
	if values.Get("issueId") == "" {
		values.Del("issueId")
	}
	pullRequest, err := client.CreatePullRequest(projectId, repositoryId, values)
	if err != nil {
		return err
	}

	fmt.Println(pullRequest.Id)

	return nil
}

func UpdatePullRequestCommand(args []string) error {
	users, err := client.GetUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Println(user.Name)
	}
	return nil
	if len(args) < 2 {
		return nil
	}

	values, err := parseMarkdownForPR(args[1])
	if err != nil {
		return err
	}
	projectId := values.Get("projectId")
	repositoryId := values.Get("repositoryId")
	values.Del("projectId")
	values.Del("repositoryId")
	number, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	pullRequest, err := client.UpdatePullRequest(projectId, repositoryId, number, values)
	if err != nil {
		return err
	}

	fmt.Println(pullRequest.Id)

	return nil
}

func DeleteIssueCommand(args []string) error {
	if len(args) < 1 {
		return nil
	}

	issueId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	issue, err := client.DeleteIssue(issueId)
	if err != nil {
		return err
	}

	fmt.Println(issue.Id)

	return nil
}

func GetPullRequestCommand(args []string) error {
	if len(args) < 2 {
		return nil
	}

	projects, err := client.GetProjects(nil)
	if err != nil {
		return nil
	}

	var projectId string
	var projectName string
	var repositoryId string

	for _, project := range projects {
		repos, err := client.GetRepositories(fmt.Sprint(project.Id), nil)
		if err != nil {
			return err
		}
		for _, repo := range repos {
			if repo.Name == args[0] {
				projectId = fmt.Sprintf("%v", project.Id)
				projectName = project.Name
				repositoryId = fmt.Sprint(repo.Id)
				break
			}
		}
	}
	if repositoryId == "" {
		return fmt.Errorf("%v not found", args[0])
	}
	number, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	pullRequest, err := client.GetPullRequest(projectId, repositoryId, number, nil)
	if err != nil {
		return err
	}

	fmt.Println("---")
	fmt.Println("summary:", pullRequest.Summary)
	fmt.Println("project:", projectName)
	fmt.Println("status:", pullRequest.Status.Name)
	fmt.Println("number:", pullRequest.Number)
	fmt.Println("branch:", pullRequest.Branch)
	fmt.Println("base:", pullRequest.Base)
	fmt.Println("issuekey:", pullRequest.Issue.IssueKey)
	fmt.Println("assignee:", pullRequest.Assignee.Name)
	fmt.Println("created:", pullRequest.CreateUser.Name)
	fmt.Println("---")
	fmt.Println(pullRequest.Description)

	return nil
}

func GetIssueCommand(args []string) error {
	if len(args) < 1 {
		return nil
	}

	issue, err := client.GetIssue(args[0])
	if err != nil {
		return err
	}

	fmt.Println("---")
	fmt.Println("summary:", issue.Summary)
	fmt.Println("parentissueid:", issue.ParentIssueId)
	fmt.Println("issuetype:", issue.IssueType.Name)
	fmt.Println("status:", issue.Status.Name)
	fmt.Println("priority:", issue.Priority.Name)
	fmt.Println("assignee:", issue.Assignee.Name)
	fmt.Println("created:", issue.CreatedUser.Name)
	fmt.Println("start:", issue.StartDate)
	fmt.Println("due:", issue.DueDate)
	fmt.Println("estimated:", issue.EstimatedHours)
	fmt.Println("actual:", issue.ActualHours)
	fmt.Printf("url: https://%s.backlog.jp/view/%s\n", spaceName, issue.IssueKey)
	fmt.Println("---")
	fmt.Println(issue.Description)

	return nil
}

func UpdateIssueCommand(args []string) error {
	if len(args) < 2 {
		return nil
	}

	values, err := parseMarkdown(args[1])
	if err != nil {
		return err
	}

	delete(values, "projectId")

	issue, err := client.SetIssue(args[0], values)
	if err != nil {
		return err
	}

	fmt.Println(issue.Id)

	return nil
}

func GetCommentsCommand(args []string) error {
	if len(args) < 1 {
		return nil
	}

	value := url.Values{}
	value.Add("count", "50")
	comments, err := client.GetComments(args[0], value)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if len(comment.ChangeLog) > 0 {
			fmt.Println(comment.CreatedUser.Name, "が課題の内容を変更しました。")
      fmt.Println(comment.Content)
			for _, change := range comment.ChangeLog {
				fmt.Println("  ", change.Field, change.OriginalValue, "->", change.NewValue)
			}
		} else {
			fmt.Println(comment.CreatedUser.Name, comment.Content)
		}
	}
	return nil
}
func VersionCommand(args []string) error {
	fmt.Printf("%v-%v\n", version, revision)

	return nil
}

func HelpCommand(args []string) error {
	fmt.Println(`usage: backlog <command> [options]

Commands:
  p, post
    Create an issue with given markdown file.
  g, get
    Print detail of the specific issue.
  u, update
      Replace existing issue with given markdown file.
    d, delete
      Delete specific issue.
  l, list
    List projects and its issues.
  v, version
    Print version and revision info.
  h, help
    Print this message.`)

	return nil
}
