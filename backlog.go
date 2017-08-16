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

	"github.com/ericaro/frontmatter"
	"github.com/moutend/go-backlog"
)

var (
	version  = "dev"
	revision = "latest"
)

func parseMarkdown(client *backlog.Client, filename string) (url.Values, error) {
	var wg sync.WaitGroup

	priorityNameToId := make(map[string]int)
	projectNameToId := make(map[string]int)
	issueTypeNameToId := make(map[string]int)
	statusNameToId := make(map[string]int)

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
			for _, issueType := range issueTypes {
				issueTypeNameToId[issueType.Name] = issueType.Id
			}
		}
		wg.Done()
	}()

	wg.Wait()

	type frontmatterOption struct {
		Project        string `fm:"project"`
		ProjectId      string `fm:"projectid"`
		IssueType      string `fm:"issuetype"`
		IssueTypeId    string `fm:"issuetypeid"`
		Priority       string `fm:"priority"`
		PriorityId     string `fm:"priorityid"`
		Status         string `fm:"status"`
		StatusId       string `fm:"statusid"`
		ParentIssue    string `fm:"parent"`
		ParentIssueId  string `fm:"parentid"`
		EstimatedHours string `fm"estimated"`
		ActualHours    string `fm"actual"`
		Summary        string `fm:"summary"`
		Description    string `fm:"content"`
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
		values.Add("issueTypeId", strconv.Itoa(issueTypeNameToId[fo.IssueType]))
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
	if fo.EstimatedHours != "" {
		values.Add("estimatedHours", fo.EstimatedHours)
	}
	if fo.ActualHours != "" {
		values.Add("actualHours", fo.ActualHours)
	}

	values.Add("summary", fo.Summary)
	values.Add("description", fo.Description)

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

	var client *backlog.Client
	var err error
	var debugFlag bool

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])
	command := args[1]
	args = f.Args()

	if client, err = backlog.New(os.Getenv("BACKLOG_SPACE"), os.Getenv("BACKLOG_TOKEN")); err != nil {
		return err
	}
	if debugFlag {
		client.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	switch command {
	case "l":
		return ListCommand(client, args)
	case "list":
		return ListCommand(client, args)
	case "p":
		return CreateIssueCommand(client, args)
	case "post":
		return CreateIssueCommand(client, args)
	case "c":
		return CreateIssueCommand(client, args)
	case "create":
		return CreateIssueCommand(client, args)
	case "d":
		return DeleteIssueCommand(client, args)
	case "delete":
		return DeleteIssueCommand(client, args)
	case "u":
		return UpdateIssueCommand(client, args)
	case "update":
		return UpdateIssueCommand(client, args)
	default:
		return fmt.Errorf("%s is not a subcommand", command)
	}
}

func ListCommand(client *backlog.Client, args []string) error {
	projects, err := client.GetProjects(nil)
	if err != nil {
		return nil
	}
	for _, project := range projects {
		fmt.Printf("- %v (id:%v)\n", project.Name, project.Id)

		query := url.Values{}
		query.Add("projectId[]", fmt.Sprintf("%v", project.Id))
		query.Add("count", "100")

		issues, err := client.GetIssues(query)
		if err != nil {
			return err
		}
		for _, issue := range issues {
			fmt.Printf("  - [%v %v] %v by @%v (id:%v)\n", issue.Status.Name, issue.IssueType.Name, issue.Summary, issue.CreatedUser.Name, issue.Id)
		}
	}

	return nil
}

func CreateIssueCommand(client *backlog.Client, args []string) error {
	if len(args) < 1 {
		return nil
	}

	values, err := parseMarkdown(client, args[0])
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

func DeleteIssueCommand(client *backlog.Client, args []string) error {
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

func UpdateIssueCommand(client *backlog.Client, args []string) error {
	if len(args) < 2 {
		return nil
	}

	values, err := parseMarkdown(client, args[1])
	if err != nil {
		return err
	}

	delete(values, "projectId")

	issueId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	issue, err := client.SetIssue(issueId, values)
	if err != nil {
		return err
	}

	fmt.Println(issue.Id)

	return nil
}

func VersionCommand(args []string) error {
	fmt.Printf("%v-%v\n", version, revision)

	return nil
}

func HelpCommand(args []string) error {
	fmt.Println(`usage: backlog <command> [options]

Commands:
  c, create
    Create an issue with given markdown file.
  p, post
    Alias of 'create' command.
  u, update
      Replace existing issue with given markdown file.
    d, delete
      Delete issue by ID.
  l, list
    List projects and its issues.
  v, version
    Print version and revision info.
  h, help
    Print this message.`)

	return nil
}
