#!/bin/bash
export GOARCH=arm
export GOARM=6

UPARG='-update=yes'
if [ ! -z "$1" ]
then
  diskutil unmountDisk "$1"
  UPARG="-overwrite=$1"
fi

~/go/bin/gokr-packer \
  -tls=self-signed \
  -kernel_package=github.com/gokrazy-community/kernel-rpi-os-32/dist \
  -hostname=streamdeckpi \
  -serial_console=serial0,115200 \
  "$UPARG" \
  github.com/gokrazy/wifi \
  github.com/Doridian/streamdeckpi/agent/cmd/agent
