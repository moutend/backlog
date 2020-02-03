package app

import (
	"os"

	"backlog/internal/backlog"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:               "backlog",
	PersistentPreRunE: rootPersistentPreRunE,
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) error {
	space := os.Getenv("BACKLOG_SPACE")
	token := os.Getenv("BACKLOG_TOKEN")

	if err := backlog.Setup(space, token); err != nil {
		return err
	}

	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		backlog.SetHTTPClient(backlog.NewDebugClient())
	}

	return nil
}

func init() {
	RootCommand.Flags().BoolP("debug", "d", false, "Enable debug output")
}
