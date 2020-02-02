package app

import (
	"backlog/internal/backlog"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/moutend/go-backlog/pkg/types"
)

func TestIssueUnmarshal(t *testing.T) {
	space := os.Getenv("BACKLOG_SPACE")
	token := os.Getenv("BACKLOG_TOKEN")

	backlog.Setup(space, token)
	backlog.SetDebug(true)

	data, err := ioutil.ReadFile("/tmp/sdc/sdc-212.md")

	if err != nil {
		t.Fatal(err)
	}

	i := &Issue{}

	if err := i.Unmarshal(data); err != nil {
		t.Fatal(err)
	}

	output, err := json.Marshal(i)

	if err != nil {
		t.Fatal(err)
	}

	log.Printf("JSON: %s\n", output)
}

func TestIssueMarshal(t *testing.T) {
	space := os.Getenv("BACKLOG_SPACE")
	token := os.Getenv("BACKLOG_TOKEN")

	backlog.Setup(space, token)
	backlog.SetDebug(true)

	i := &Issue{
		Issue: &types.Issue{
			Summary:        "Issue Summary",
			Description:    "Issue description",
			StartDate:      types.NewDate("2020-02-02"),
			EstimatedHours: types.NewHours(1.5),
		},
	}

	output, err := i.Marshal()

	if err != nil {
		t.Fatal(err)
	}

	log.Printf("markdown: %s\n", output)
}
