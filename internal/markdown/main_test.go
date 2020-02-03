package markdown

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(0)

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
