#!/bin/sh
set -ex
go install github.com/gokrazy/tools/cmd/gok@main
cd agent && rm -fv _gokrazy/extrafiles/etc/streamdeckpi/*.yml && ./templategen.py && cd ..
rm -rf builddir
~/go/bin/gok update --instance streamdeckpi
