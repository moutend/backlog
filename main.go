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
	if issues, err = client.Issues(query); err != nil {
		return
	}

	JST, _ := time.LoadLocation("Asia/Tokyo")
	for _, issue := range issues {
		fmt.Println(issue.Summary, issue.Created.Time().In(JST))
	}

	return
}
