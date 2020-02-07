package app

import (
	"github.com/spf13/cobra"
)

const (
	version = "v0.1.0"
	commit  = "latest"
)

var versionCommand = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	RunE:    versionCommandRunE,
}

func versionCommandRunE(cmd *cobra.Command, args []string) error {
	cmd.Printf("%s-%s\n", version, commit)
	return nil
}

func init() {
	RootCommand.AddCommand(versionCommand)
}
