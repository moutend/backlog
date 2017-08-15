# backlog

[![GitHub release](https://img.shields.io/github/release/moutend/backlog.svg?style=flat-square)][release]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![CircleCI](https://circleci.com/gh/moutend/backlog.svg?style=svg&circle-token=e7748578056ded93a5532904c047fc0f23db3bba)](https://circleci.com/gh/moutend/backlog)

[release]: https://github.com/moutend/backlog/releases
[license]: https://github.com/moutend/backlog/blob/master/LICENSE
[status]: https://circleci.com/gh/moutend/backlog

`backlog` is a CLI client for [https://backlog.jp](https://backlog.jp).

# Usage

To list projects and its issues, use `list` command:

```console
backlog list
```

To create an issue, write a markdown file like the following at first.

**issue.md**

```markdown
---
summary: TODO
project: My project
issuetype: task
priority: low
status: ongoing
---

# TODO

- Implement XXX
- Fix YYY
- Delete ZZZ
```

And then use `create` command to post the issue.

```console
backlog create issue.md
```

For more information, see `help` command.

# Installation

## Windows / Linux

You can download the executable for 32 bit / 64 bit at [GitHub releases page](https://github.com/moutend/backlog/releases/).

## Mac

The easiest way is Homebrew.

```shell
$ brew tap moutend/homebrew-backlog
$ brew install backlog
```

## `go install`

If you have already set up Go environment, just `go install`.

```shell
$ go install github.com/moutend/backlog
```

## Author

[Yoshiyuki Koyanagi](https://github.com/moutend)
