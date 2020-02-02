package testutil

import "os"

var (
	BacklogSpace       string
	BacklogToken       string
	EnableHTTPRequest  bool
	EnableLoggerOutput bool
)

func init() {
	BacklogSpace = os.Getenv("TEST_BACKLOG_SPACE")
	BacklogToken = os.Getenv("TEST_BACKLOG_TOKEN")
	EnableHTTPRequest = os.Getenv("TEST_ENABLE_HTTP_REQUEST") != ""
	EnableLoggerOutput = os.Getenv("TEST_ENABLE_LOGGER_OUTPUT") != ""

	if BacklogSpace == "" {
		BacklogSpace = "backlog-space"
	}
	if BacklogToken == "" {
		BacklogToken = "backlog-token"
	}
}
