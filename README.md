# Pimoroni Fan SHIM controller written in Golang

![Go](https://github.com/jsiebens/fanshim-go/workflows/Go/badge.svg)

A small application with temperature monitoring and fan control

Supported arguments:

* `--on-threshold N` the temperature at which to turn the fan on, in degrees C (default 65)
* `--off-threshold N` the temperature at which to turn the fan off, in degrees C (default 55)
* `--delay N` the delay between subsequent temperature readings, in seconds (default 2)

## Alternate Software

* Pimoroni Fan SHIM in Python - https://github.com/pimoroni/fanshim-python
* Fan SHIM in C, using WiringPi - https://github.com/flobernd/raspi-fanshim
* Fan SHIM in C++, using libgpiod - https://github.com/daviehh/fanshim-cpp
