package app

import (
	"net/url"
	"sort"
	"strconv"

	"github.com/moutend/backlog/internal/backlog"
	"github.com/moutend/backlog/internal/cache"
	"github.com/moutend/go-backlog/pkg/types"
	"github.com/spf13/cobra"
)

var commentCommand = &cobra.Command{
	Use:     "comment",
	Aliases: []string{"c"},
	RunE:    commentCommandRunE,
}

func commentCommandRunE(cmd *cobra.Command, args []string) error {
	return nil
}

var commentIssueCommand = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"i"},
	RunE:    commentIssueCommandRunE,
}

func commentIssueCommandRunE(cmd *cobra.Command, args []string) error {
	var (
		comments []*types.Comment
		err      error
	)

	timeout, _ := cmd.Flags().GetDuration("timeout")

	if timeout == 0 {
		goto PRINT_ISSUE_COMMENTS
	}

	{
		query := url.Values{}

		query.Add("count", "100")

		comments, err = backlog.GetIssueComments(args[0], query)

		if err != nil {
			warn.Println(err)

			goto PRINT_ISSUE_COMMENTS
		}
		if err := cache.SaveIssueComments(args[0], comments); err != nil {
			return err
		}
	}

PRINT_ISSUE_COMMENTS:

	comments, err = cache.LoadIssueComments(args[0])

	if err != nil {
		return err
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Created.Time().After(comments[j].Created.Time())
	})

	for _, comment := range comments {
		cmd.Printf(
			"%s %s",
			comment.Created,
			comment.CreatedUser.Name,
		)
		if comment.Content != "" {
			cmd.Printf(" %s\n", comment.Content)
		} else {
			cmd.Printf("\n")
		}
		for _, changeLog := range comment.ChangeLog {
			if changeLog.Field != "" {
				cmd.Printf("\n%s\n", changeLog.Field)
			}
			if changeLog.NewValue != nil {
				cmd.Println("+", *changeLog.NewValue)
			}
			if changeLog.OriginalValue != nil {
				cmd.Println("-", *changeLog.OriginalValue)
			}
			if changeLog.NotificationInfo != nil {
				cmd.Println(changeLog.NotificationInfo.Type)
			}
		}
	}

	return nil
}

var commentPullRequestCommand = &cobra.Command{
	Use:     "pullrequest",
	Aliases: []string{"pr"},
	RunE:    commentPullRequestCommandRunE,
}

func commentPullRequestCommandRunE(cmd *cobra.Command, args []string) error {
	var (
		comments []*types.Comment
	)

	if len(args) < 3 {
		return nil
	}

	projectKey := args[0]
	repositoryName := args[1]
	number, err := strconv.ParseInt(args[2], 10, 64)

	if err != nil {
		return err
	}

	timeout, _ := cmd.Flags().GetDuration("timeout")

	if timeout == 0 {
		goto PRINT_PULLREQUEST_COMMENTS
	}
	{
		query := url.Values{}

		query.Add("count", "100")

		comments, err = backlog.GetPullRequestComments(projectKey, repositoryName, number, query)

		if err != nil {
			warn.Println(err)

			goto PRINT_PULLREQUEST_COMMENTS
		}
		if err := cache.SavePullRequestComments(projectKey, repositoryName, number, comments); err != nil {
			return err
		}
	}

PRINT_PULLREQUEST_COMMENTS:

	comments, err = cache.LoadPullRequestComments(projectKey, repositoryName, number)

	if err != nil {
		return err
	}
	for _, comment := range comments {
		cmd.Printf(
			"%s (%s)\n",
			comment.CreatedUser.Name,
			comment.Created,
		)
		for _, changeLog := range comment.ChangeLog {
			cmd.Println(changeLog.NewValue)

			if changeLog.AttachmentInfo != nil {
				cmd.Println(changeLog.AttachmentInfo)
			}
		}
	}
	return nil
}

func init() {
	RootCommand.AddCommand(commentCommand)

	commentCommand.AddCommand(commentIssueCommand)
	commentCommand.AddCommand(commentPullRequestCommand)
}
