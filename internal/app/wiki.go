package app

import (
	"backlog/internal/markdown"
	"fmt"
	"strconv"

	"github.com/moutend/go-backlog/pkg/types"

	"backlog/internal/backlog"

	"github.com/spf13/cobra"
)

var wikiCommand = &cobra.Command{
	Use:     "wiki",
	Aliases: []string{"w"},
	RunE:    wikiCommandRunE,
}

func wikiCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var wikiListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    wikiListCommandRunE,
}

func wikiListCommandRunE(cmd *cobra.Command, args []string) error {
	projects, err := backlog.GetProjects(nil)

	if err != nil {
		return err
	}

	wikiMap := map[uint64][]*types.Wiki{}

	for _, project := range projects {
		wikis, err := backlog.GetWikis(project.ProjectKey, nil)

		if err != nil {
			return err
		}

		wikiMap[project.Id] = wikis
	}
	for i, project := range projects {
		cmd.Printf("# [%s] %s\n\n", project.ProjectKey, project.Name)

		wikis := wikiMap[project.Id]

		if len(wikis) == 0 {
			cmd.Println("No wikis.")

			goto NEXT
		}
		for _, wiki := range wikis {
			cmd.Printf("- %s (id:%d)\n", wiki.Name, wiki.Id)
		}

	NEXT:

		if i < len(projects)-1 {
			cmd.Println("")
		}
	}

	return nil
}

var wikiReadCommand = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	RunE:    wikiReadCommandRunE,
}

func wikiReadCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	i, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	wiki, err := backlog.GetWiki(uint64(i))

	if err != nil {
		return err
	}

	project, err := backlog.GetProject(fmt.Sprint(wiki.ProjectId))

	if err != nil {
		return err
	}

	mw := markdown.Wiki{
		Wiki:    wiki,
		Project: project,
	}

	output, err := mw.Marshal()

	if err != nil {
		return err
	}

	cmd.Printf("%s", output)

	return nil
}

func init() {
	RootCommand.AddCommand(wikiCommand)

	wikiCommand.AddCommand(wikiListCommand)
	wikiCommand.AddCommand(wikiReadCommand)
}
