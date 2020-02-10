package app

import (
	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/backlog/internal/cache"
	"github.com/moutend/go-backlog/pkg/types"
	"github.com/spf13/cobra"
)

var myselfCommand = &cobra.Command{
	Use:     "myself",
	Aliases: []string{"m"},
	RunE:    myselfCommandRunE,
}

func myselfCommandRunE(cmd *cobra.Command, args []string) error {
	var (
		myself *types.User
		err    error
	)

	myself, err = backlog.GetMyself()

	if err != nil {
		goto PRINT_MYSELF
	}
	if err := cache.SaveMyself(myself); err != nil {
		return err
	}

PRINT_MYSELF:

	myself, err = cache.LoadMyself()

	if err != nil {
		return err
	}

	cmd.Printf("%s\n", myself.Name)

	return nil
}

func init() {
	RootCommand.AddCommand(myselfCommand)
}
