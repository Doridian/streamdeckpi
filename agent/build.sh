#!/bin/sh

export GOARM=6
export GOARCH=arm
export GOOS=linux
#export CGO_ENABLED=1
#export CC=/opt/homebrew/bin/arm-none-eabi-gcc

gok run --instance streamdeckpi
