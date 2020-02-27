package backlog

import (
	"net/http"

	"github.com/moutend/go-backlog/pkg/client"
)

var (
	backlogClient    *client.Client
	backlogSpaceName string
)

func Setup(space, token string, debugClient *http.Client) error {
	c, err := client.New(space, token, client.OptionHTTPClient(debugClient))

	if err != nil {
		return err
	}

	backlogSpaceName = space
	backlogClient = c

	return nil
}

func SpaceName() string {
	return backlogSpaceName
}
