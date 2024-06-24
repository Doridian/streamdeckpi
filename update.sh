#!/bin/sh
set -ex
go install github.com/gokrazy/tools/cmd/gok@main
rm -rf builddir
~/go/bin/gok update --instance streamdeckpi
