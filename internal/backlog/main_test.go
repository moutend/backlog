package backlog

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"backlog/internal/testutil"
)

func TestMain(m *testing.M) {
	log.SetPrefix("TEST: ")
	log.SetFlags(0)

	if !testutil.EnableLoggerOutput {
		log.SetOutput(ioutil.Discard)
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
