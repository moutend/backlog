package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/backlog/internal/cache"
	"github.com/moutend/backlog/internal/markdown"
	"github.com/moutend/go-backlog/pkg/types"
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

	sortedBy, _ := cmd.Flags().GetString("sort")

	switch strings.ToLower(sortedBy) {
	case "created":
		sort.Slice(wikis, func(i, j int) bool {
			return wikis[i].Created.Time().After(wikis[j].Created.Time())
		})
	case "updated":
		sort.Slice(wikis, func(i, j int) bool {
			return wikis[i].Updated.Time().After(wikis[j].Updated.Time())
		})
	default:
		sort.Slice(wikis, func(i, j int) bool {
			return wikis[i].Id > wikis[j].Id
		})
	}
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

var wikiCreateCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	RunE:    wikiCreateCommandRunE,
}

func wikiCreateCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	data, err := ioutil.ReadFile(args[0])

	if err != nil {
		return err
	}

	mw := &markdown.Wiki{}

	if err := mw.Unmarshal(data); err != nil {
		return err
	}

	createdWiki, err := backlog.AddWiki(mw.Wiki, true)

	if err != nil {
		return err
	}

	cmd.Println("Created wiki:", createdWiki.Id)

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

var wikiUpdateCommand = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	RunE:    wikiUpdateCommandRunE,
}

func wikiUpdateCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	data, err := ioutil.ReadFile(args[0])

	if err != nil {
		return err
	}

	mw := &markdown.Wiki{}

	if err := mw.Unmarshal(data); err != nil {
		return err
	}

	updatedWiki, err := backlog.UpdateWiki(mw.Wiki, true)

	if err != nil {
		return err
	}

	cmd.Println("Updated wiki:", updatedWiki.Id)

	return nil
}

var wikiDeleteCommand = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	RunE:    wikiDeleteCommandRunE,
}

func wikiDeleteCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	i, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	deletedWiki, err := backlog.DeleteWiki(uint64(i))

	if err != nil {
		return err
	}

	cmd.Println("Deleted wiki:", deletedWiki.Id)

	return nil
}

var wikiAttachmentCommand = &cobra.Command{
	Use:     "attachment",
	Aliases: []string{"a"},
	RunE:    wikiAttachmentCommandRunE,
}

func wikiAttachmentCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var wikiAttachmentListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE:    wikiAttachmentListCommandRunE,
}

func wikiAttachmentListCommandRunE(cmd *cobra.Command, args []string) error {
	var (
		wiki        *types.Wiki
		attachments []*types.Attachment
		err         error
	)
	if len(args) < 0 {
		return nil
	}

	i, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	wiki, err = backlog.GetWiki(uint64(i))

	if err != nil {
		goto PRINT_ATTACHMENTS
	}
	if err := cache.Save(wiki); err != nil {
		return err
	}

	attachments, err = backlog.GetWikiAttachments(wiki.Id)

	if err != nil {
		goto PRINT_ATTACHMENTS
	}
	if err := cache.SaveWikiAttachments(wiki, attachments); err != nil {
		return err
	}

PRINT_ATTACHMENTS:

	wiki, err = cache.LoadWiki(wiki.Id)

	if err != nil {
		return err
	}

	attachments, err = cache.LoadWikiAttachments(wiki.Id)

	if err != nil {
		return err
	}

	cmd.Printf("# %s (%d)\n", wiki.Name, wiki.Id)

	if len(attachments) == 0 {
		return nil
	}

	cmd.Printf("\n")

	for _, attachment := range attachments {
		cmd.Printf("- %s (id:%d)\n", attachment.Name, attachment.Id)
	}

	return nil
}

var wikiAttachmentCreateCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	RunE:    wikiAttachmentCreateCommandRunE,
}

func wikiAttachmentCreateCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
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

	attachments, err := backlog.AddWikiAttachments(wiki.Id, args[1:]...)

	if err != nil {
		return err
	}

	for _, attachment := range attachments {
		cmd.Printf("Created wiki attachment: %s (id:%d)\n", attachment.Name, attachment.Id)
	}

	return nil
}

var wikiAttachmentReadCommand = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	RunE:    wikiAttachmentReadCommandRunE,
}

func wikiAttachmentReadCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
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

	i, err = strconv.Atoi(args[1])

	if err != nil {
		return err
	}

	data, filename, err := backlog.DownloadWikiAttachment(wiki.Id, uint64(i))

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	cmd.Println("Downloaded wiki attachment:", filename)

	return nil
}

var wikiAttachmentDeleteCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	RunE:    wikiAttachmentDeleteCommandRunE,
}

func wikiAttachmentDeleteCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	wikiAttachmentCommand.AddCommand(wikiAttachmentListCommand)
	wikiAttachmentCommand.AddCommand(wikiAttachmentCreateCommand)
	wikiAttachmentCommand.AddCommand(wikiAttachmentReadCommand)
	wikiAttachmentCommand.AddCommand(wikiAttachmentDeleteCommand)
	wikiCommand.AddCommand(wikiAttachmentCommand)

	wikiListCommand.Flags().BoolP("desc", "", true, "Print wikis descending order")
	wikiListCommand.Flags().BoolP("asc", "", false, "Print wikis ascending order")
	wikiListCommand.Flags().StringP("sort", "", "", "Specify sorting order")

	wikiCommand.AddCommand(wikiListCommand)
	wikiCommand.AddCommand(wikiCreateCommand)
	wikiCommand.AddCommand(wikiReadCommand)
	wikiCommand.AddCommand(wikiUpdateCommand)
	wikiCommand.AddCommand(wikiDeleteCommand)

	RootCommand.AddCommand(wikiCommand)
}
