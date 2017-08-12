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
		Project     string `fm:"project"`
		ProjectId   string `fm:"projectid"`
		IssueType   string `fm:"issuetype"`
		IssueTypeId string `fm:"issuetypeid"`
		Priority    string `fm:"priority"`
		PriorityId  string `fm:"priorityid"`
		Status      string `fm:"status"`
		StatusId    string `fm:"statusid"`

		ParentIssue   string `fm:"parent"`
		ParentIssueId string `fm:"parentid"`
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
		values.Add("projectId", strconv.Itoa(projectNameToId[fo.Project]))
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

func run(args []string) (err error) {
	if len(args) < 2 {
		return HelpCommand(args)
	}
	switch args[1] {
	case "v":
		err = VersionCommand(args)
	case "version":
		err = VersionCommand(args)
	case "h":
		err = HelpCommand(args)
	case "help":
		err = HelpCommand(args)
	case "l":
		err = ListCommand(args)
	case "list":
		err = ListCommand(args)
	case "c":
		err = CreateIssueCommand(args)
	case "create":
		err = CreateIssueCommand(args)
	case "u":
		err = UpdateIssueCommand(args)
	case "update":
		err = UpdateIssueCommand(args)
	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s subcommand.\n", args[0], args[1], args[0])
	}
	return
}

func ListCommand(args []string) (err error) {
	var debugFlag bool

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])
	args = f.Args()

	var client *backlog.Client
	if client, err = backlog.New(os.Getenv("BACKLOG_SPACE"), os.Getenv("BACKLOG_TOKEN")); err != nil {
		return
	}
	if debugFlag {
		client.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}

	var projects []*backlog.Project
	if projects, err = client.GetProjects(nil); err != nil {
		return
	}
	for _, project := range projects {
		fmt.Printf("- %v (id:%v)\n", project.Name, project.Id)

		query := url.Values{}
		query.Add("projectId[]", fmt.Sprintf("%v", project.Id))
		query.Add("count", "100")

		var issues []*backlog.Issue
		if issues, err = client.GetIssues(query); err != nil {
			return
		}
		for _, issue := range issues {
			fmt.Printf("  - [%v %v] %v by @%v (id:%v)\n", issue.Status.Name, issue.IssueType.Name, issue.Summary, issue.CreatedUser.Name, issue.Id)
		}
	}

	return
}

func CreateIssueCommand(args []string) error {
	var debugFlag bool

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])
	args = f.Args()

	client, err := backlog.New(os.Getenv("BACKLOG_SPACE"), os.Getenv("BACKLOG_TOKEN"))
	if err != nil {
		return err
	}
	if debugFlag {
		client.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
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

func UpdateIssueCommand(args []string) error {
	var debugFlag bool
	var issueIdFlag int

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.IntVar(&issueIdFlag, "i", 0, "specify issue ID")
	f.Parse(args[2:])
	args = f.Args()

	client, err := backlog.New(os.Getenv("BACKLOG_SPACE"), os.Getenv("BACKLOG_TOKEN"))
	if err != nil {
		return err
	}
	if debugFlag {
		client.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	if len(args) < 1 {
		return nil
	}

	values, err := parseMarkdown(client, args[0])
	if err != nil {
		return err
	}
	delete(values, "projectId")
	issue, err := client.SetIssue(issueIdFlag, values)
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
	return nil
}
