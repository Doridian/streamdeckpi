#!/bin/sh
set -e

echo 'Building...'
rm -f ./agent
go build ./cmd/agent
echo 'Running...'
STREAMDECK_CONFIG_DIR=./_gokrazy/extrafiles/etc/streamdeckpi ./agent
