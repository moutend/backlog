package backlog

import (
	"log"
	"testing"

	"backlog/internal/testutil"
)

func TestSetup(t *testing.T) {
	if err := Setup(testutil.BacklogSpace, testutil.BacklogToken); err != nil {
		t.Fatal(err)
	}

	SetHTTPClient(testutil.NewTestClient([]byte(`{}`)))

	is, err := GetIssues(nil)

	if err != nil {
		t.Fatal(err)
	}

	for _, i := range is {
		log.Printf("GetIssues: %+v\n", i)
	}
}
