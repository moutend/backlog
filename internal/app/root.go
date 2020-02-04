package app

import (
	"os"
	"time"

	"backlog/internal/backlog"
	"backlog/internal/cache"

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
	if err := cache.Setup(space); err != nil {
		return err
	}
	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		backlog.SetHTTPClient(backlog.NewDebugClient())
	}

	return nil
}

func init() {
	RootCommand.PersistentFlags().BoolP("debug", "d", false, "Enable debug output")
	RootCommand.PersistentFlags().DurationP("timeout", "t", time.Minute, "Set timeout value (default=60s)")
}
