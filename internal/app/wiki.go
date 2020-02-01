package app

import (
	"backlog/internal/cache"
	"backlog/internal/markdown"
	"context"
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
	var (
		projects []*types.Project
		wikis    []*types.Wiki
		ctx      context.Context
		err      error
	)

	timeout, _ := cmd.Flags().GetDuration("timeout")

	if timeout == 0 {
		goto PRINT_WIKIS
	}

	ctx, _ = context.WithTimeout(context.Background(), timeout)

	projects, err = backlog.GetProjects(nil)

	if err != nil {
		warn.Println(err)

		goto PRINT_WIKIS
	}
	if err := cache.Save(projects); err != nil {
		return err
	}

	for _, project := range projects {
		wikis, err = backlog.GetWikisContext(ctx, project.ProjectKey, nil)

		if err != nil {
			warn.Println(err)

			goto PRINT_WIKIS
		}
		if err := cache.Save(wikis); err != nil {
			return err
		}
	}

PRINT_WIKIS:

	projects, err = cache.LoadProjects()

	if err != nil {
		return err
	}

	wikis, err = cache.LoadWikis()

	if err != nil {
		return err
	}

	wikisMap := map[uint64][]*types.Wiki{}

	for _, wiki := range wikis {
		wikisMap[wiki.ProjectId] = append(wikisMap[wiki.ProjectId], wiki)
	}
	for i, project := range projects {
		cmd.Printf("# [%s] %s\n\n", project.ProjectKey, project.Name)

		wikis := wikisMap[project.Id]

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
