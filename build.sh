#!/bin/sh
export GOARCH=arm
export GOARM=6

#  -overwrite /dev/disk4 \

~/go/bin/gokr-packer \
  -tls=self-signed \
  -kernel_package=github.com/Doridian/kernel-rpi-os-32/dist \
  -hostname streamdeckpi \
  -serial_console serial0,115200 \
  -update=yes \
  github.com/gokrazy/hello \
  github.com/gokrazy/breakglass \
  github.com/gokrazy/serial-busybox \
  github.com/gokrazy/wifi
