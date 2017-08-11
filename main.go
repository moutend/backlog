package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/ericaro/frontmatter"
	"github.com/moutend/go-backlog"
)

func parseMarkdown(filename string) (values url.Values, err error) {
	type FrontmatterOption struct {
		ProjectId     string `fm:"projectId"`
		IssueTypeId   string `fm:"issueTypeId"`
		PriorityId    string `fm:"priorityId"`
		ParentIssueId string `fm:"parentIssueId"`
		Summary       string `fm:"summary"`
		Description   string `fm:"content"`
	}

	var file []byte
	if file, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	fo := &FrontmatterOption{}
	if err = frontmatter.Unmarshal(file, fo); err != nil {
		return
	}
	values = url.Values{}
	values.Add("projectId", fo.ProjectId)
	values.Add("issueTypeId", fo.IssueTypeId)
	values.Add("priorityId", fo.PriorityId)
	values.Add("summary", fo.Summary)
	values.Add("description", fo.Description)

	if fo.ParentIssueId != "" {
		values.Add("parentIssueId", fo.ParentIssueId)
	}
	return
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
	case "p":
		err = PostCommand(args)
	case "post":
		err = PostCommand(args)
	case "l":
		err = ListCommand(args)
	case "list":
		err = ListCommand(args)
	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s subcommand.\n", args[0], args[1], args[0])
	}
	return
}

func ListCommand(args []string) (err error) {
	var client *backlog.Client
	if client, err = backlog.New(os.Getenv("BACKLOG_SPACE"), os.Getenv("BACKLOG_TOKEN")); err != nil {
		return
	}

	//client.SetLogger(log.New(os.Stdout, "debug: ", 0))

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

func PostCommand(args []string) (err error) {
	var debugFlag bool

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])
	args = f.Args()

	if len(args) < 1 {
		return
	}

	var values url.Values
	if values, err = parseMarkdown(args[0]); err != nil {
		return
	}

	var client *backlog.Client
	if client, err = backlog.New(os.Getenv("BACKLOG_SPACE"), os.Getenv("BACKLOG_TOKEN")); err != nil {
		return
	}
	if debugFlag {
		client.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}

	var issue *backlog.Issue
	if issue, err = client.CreateIssue(values); err != nil {
		return
	}

	fmt.Printf("post: [%v %v] %v by @%v (id:%v)\n", issue.Status.Name, issue.IssueType.Name, issue.Summary, issue.CreatedUser.Name, issue.Id)

	return
}

func VersionCommand(args []string) (err error) {
	return
}

func HelpCommand(args []string) (err error) {
	return
}
