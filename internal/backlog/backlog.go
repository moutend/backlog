package backlog

import (
	"github.com/moutend/go-backlog/pkg/client"
)

var (
	bc *client.Client
)

func New(space, token string) error {
	c, err := client.New(space, token)

	if err != nil {
		return err
	}

	bc = c

	return nil
}
