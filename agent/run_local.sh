#!/bin/sh
set -e

echo 'Building...'
rm -f ./agent
go build ./cmd/agent
echo 'Running...'
./agent
