package app

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/backlog/internal/cache"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:               "backlog",
	PersistentPreRunE: rootPersistentPreRunE,
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) error {
	switch cmd.Name() {
	case "v", "version":
		return nil
	}

	space := os.Getenv("BACKLOG_SPACE")
	token := os.Getenv("BACKLOG_TOKEN")
	debugClient := &http.Client{}

	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		debugClient = backlog.NewDebugClient()
	}
	if err := backlog.Setup(space, token, debugClient); err != nil {
		return err
	}
	if err := cache.Setup(space); err != nil {
		return err
	}
	if yes, _ := cmd.Flags().GetBool("warn"); yes {
		warn = log.New(os.Stderr, "warn: ", 0)
	} else {
		warn = log.New(ioutil.Discard, "warn: ", 0)
	}
	if timeout, _ := cmd.Flags().GetDuration("timeout"); timeout == 0 {
		return nil
	}

	myself, err := backlog.GetMyself()

	if err != nil {
		return err
	}
	if err := cache.SaveMyself(myself); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCommand.PersistentFlags().BoolP("debug", "d", false, "Enable debug output")
	RootCommand.PersistentFlags().BoolP("warn", "w", true, "Enable warn output")
	RootCommand.PersistentFlags().DurationP("timeout", "t", 30*time.Second, "Set timeout value")
}
