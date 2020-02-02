#!/bin/bash

set -e

cat ~/github.com/moutend/go-backlog/pkg/client/*.go | ./generate.rb > methods.go
gofmt -l -w . && goimports -l -w .
