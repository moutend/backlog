# backlog

[![GitHub release](https://img.shields.io/github/release/moutend/backlog.svg?style=flat-square)][release]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![CircleCI](https://circleci.com/gh/moutend/backlog.svg?style=svg&circle-token=e7748578056ded93a5532904c047fc0f23db3bba)](https://circleci.com/gh/moutend/backlog)

[release]: https://github.com/moutend/backlog/releases
[license]: https://github.com/moutend/backlog/blob/master/LICENSE
[status]: https://circleci.com/gh/moutend/backlog

`backlog` is a CLI client for [https://backlog.jp](https://backlog.jp).

## Prerequisites

- Go 1.13 or later

## Download

Use `go get` to download the executable.

```console
go get -u github.com/moutend/backlog/cmd/backlog
```

Or build and install by hand.

```console
git clone https://github.com/moutend/backlog
cd ./backlog/cmd/backlog
go build
go install
```

## Usage

Before run the `backlog`, set the following environment variables.

- `BACKLOG_SPACE` ... Your backlog domain except `https://`.
- `BACKLOG_TOKEN` ... API token.

For example, set the environment variables like:

```console
export BACKLOG_SPACE='example.backlog.com'
export BACKLOG_TOKEN='xxxxxxxx'
```

Now, you're able to run the `backlog`.

```console
backlog project list
```

For more details, read `backlog help`.

## Author

[Yoshiyuki Koyanagi](https://github.com/moutend)

## LICENSE

MIT
