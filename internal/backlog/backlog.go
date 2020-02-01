package backlog

import (
	"github.com/moutend/go-backlog/pkg/client"
)

var (
	backlogClient    *client.Client
	backlogSpaceName string
)

func Setup(space, token string) error {
	c, err := client.New(space, token)

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

func SetHTTPClient(hc client.HTTPClient) {
	backlogClient.SetHTTPClient(hc)
}
