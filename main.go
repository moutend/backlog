package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/moutend/backlog"
)

func main() {
	var err error
	if err = run(os.Args); err != nil {
		log.Fatal(err)
	}
	return
}

func run(args []string) (err error) {
	var client *backlog.Client

	if client, err = backlog.New("abby", os.Getenv("BACKLOG_TOKEN")); err != nil {
		return
	}

	var issues []*backlog.Issue

	query := url.Values{
		"count": []string{"10"},
	}
	if issues, err = client.GetIssues(query); err != nil {
		return
	}

	issuesMap := make(map[int]*backlog.Issue)
	timezone, _ := time.LoadLocation("Asia/Tokyo")

	for _, issue := range issues {
		issuesMap[issue.Id] = issue
	}
	for _, issue := range issues {
		var parentIssue *backlog.Issue

		if issue.ParentIssueId != 0 {
			parentIssue, _ = issuesMap[issue.ParentIssueId]
		}

		fmt.Printf("# [%s] %s\n\n", issue.IssueType.Name, issue.Summary)

		if parentIssue != nil {
			fmt.Println("- 親課題: ", parentIssue.Summary, parentIssue.Id)
		}
		if issue.ParentIssueId != 0 && parentIssue == nil {
			fmt.Println("- 親課題: ", issue.ParentIssueId)
		}

		fmt.Println("- 状態:", issue.Status.Name)
		fmt.Println("- 作成者:", issue.CreatedUser.Name)
		fmt.Println("- 作成日時:", issue.Created.Time().In(timezone))
		fmt.Println("- 更新日時:", issue.Updated.Time().In(timezone))

		fmt.Println("\n", issue.Description, "\n")
	}

	return
}
