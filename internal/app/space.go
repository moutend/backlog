package app

import (
	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/backlog/internal/cache"
	"github.com/spf13/cobra"
)

var spaceCommand = &cobra.Command{
	Use:     "space",
	Aliases: []string{"s"},
	RunE:    spaceCommandRunE,
}

func spaceCommandRunE(cmd *cobra.Command, args []string) error {
	myself, err := cache.LoadMyself()

	if err != nil {
		return err
	}

	cmd.Printf("%s\n", myself.Name)

	return nil
}

var spaceDiskUsageCommand = &cobra.Command{
	Use:     "disk",
	Aliases: []string{"d"},
	RunE:    spaceDiskUsageCommandRunE,
}

func spaceDiskUsageCommandRunE(cmd *cobra.Command, args []string) error {
	totalDiskUsage, err := backlog.GetSpaceDiskUsage()

	if err != nil {
		return err
	}

	cmd.Printf("- Capacity %s\n", toHumanReadable(totalDiskUsage.Capacity))
	cmd.Printf("- Issue %s\n", toHumanReadable(totalDiskUsage.Issue))
	cmd.Printf("- Wiki %s\n", toHumanReadable(totalDiskUsage.Wiki))
	cmd.Printf("- File %s\n", toHumanReadable(totalDiskUsage.File))
	cmd.Printf("- SubVersion %s\n", toHumanReadable(totalDiskUsage.SubVersion))
	cmd.Printf("- Git %s\n", toHumanReadable(totalDiskUsage.Git))
	cmd.Printf("- GitLFS %s\n", toHumanReadable(totalDiskUsage.GitLFS))

	return nil
}

func init() {
	spaceCommand.AddCommand(spaceDiskUsageCommand)

	RootCommand.AddCommand(spaceCommand)
}
