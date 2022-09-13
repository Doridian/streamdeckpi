#!/bin/sh
export GOARCH=arm
export GOARM=6

#  -overwrite /dev/disk4 \

~/go/bin/gokr-packer \
  -tls=self-signed \
  -kernel_package=github.com/gokrazy-community/kernel-rpi-os-32/dist \
  -hostname streamdeckpi \
  -serial_console serial0,115200 \
  -update=yes \
  github.com/gokrazy/wifi \
  github.com/Doridian/streamdeckpi/agent/cmd/agent
