#!/bin/sh
set -e

echo 'Building...'
./templategen.py
rm -f ./agent
go build -o ./agent ./
echo 'Running...'
STREAMDECKPI_CONFIG_DIR=./_gokrazy/extrafiles/etc/streamdeckpi ./agent
