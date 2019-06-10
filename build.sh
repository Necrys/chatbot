#!/bin/sh 

go build -ldflags "-s \
  -X main.Version=1.0 \
  -X main.Commit=`git rev-list --all --count` \
  -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%S'` \
  -X main.GitHash=`git rev-parse HEAD`"
