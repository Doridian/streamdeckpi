#!/bin/sh

export GOARM=6
export GOARCH=arm
export GOOS=linux
export CGO_ENABLED=0

gok run --instance streamdeckpi --project ./cmd/agent
