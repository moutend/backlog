package backlog

import (
	"log"
	"net/http"
	"os"

	"github.com/moutend/go-backlog/pkg/client"
)

var (
	bc               *client.Client
	backlogSpaceName string
)

func Setup(space, token string) error {
	c, err := client.New(space, token)

	if err != nil {
		return err
	}

	backlogSpaceName = space
	bc = c

	return nil
}

func SpaceName() string {
	return backlogSpaceName
}

func SetDebug(debug bool) {
	if debug {
		hc := &HTTPClient{
			logger:     log.New(os.Stdout, "DEBUG: ", 0),
			httpClient: &http.Client{},
		}

		bc.SetHTTPClient(hc)
	} else {
		hc := &http.Client{}

		bc.SetHTTPClient(hc)
	}
}
