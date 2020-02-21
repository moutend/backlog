package app

import (
	"github.com/moutend/backlog/internal/cache"
	"github.com/spf13/cobra"
)

var myselfCommand = &cobra.Command{
	Use:     "myself",
	Aliases: []string{"m"},
	RunE:    myselfCommandRunE,
}

func myselfCommandRunE(cmd *cobra.Command, args []string) error {
	myself, err := cache.LoadMyself()

	if err != nil {
		return err
	}

	cmd.Printf("%s\n", myself.Name)

	return nil
}

func init() {
	RootCommand.AddCommand(myselfCommand)
}
