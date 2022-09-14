#!/bin/sh

export GOARM=6
export GOARCH=arm
export GOOS=linux

cd ./cmd/agent && gok run --instance streamdeckpi

unset GOARM
unset GOARCH
unset GOOS
