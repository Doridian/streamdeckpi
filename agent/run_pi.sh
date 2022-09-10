#!/bin/sh

export GOARM=6
export GOARCH=arm
export GOOS=linux

gok run --instance streamdeckpi --project ./cmd/agent

unset GOARM
unset GOARCH
unset GOOS
