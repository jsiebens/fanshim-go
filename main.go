package main

import (
	"flag"
	"fmt"
	"github.com/ikester/gpio"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/shirou/gopsutil/host"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const FanControl = 18
const LedData int = 15
const LedClock int = 14

var offThreshold = flag.Float64("off-threshold", 55.0, "Temperature threshold in degrees C to disable fan")
var onThreshold = flag.Float64("on-threshold", 65.0, "Temperature threshold in degrees C to enable fan")
var extendedColors = flag.Bool("extended-colors", false, "Extend LED colors for outside of normal low to high range")
var delay = flag.Int("delay", 2, "Delay, in seconds, between temperature readings")
var verbose = flag.Bool("verbose", false, "Output temp and fan status messages")

func main() {
	flag.Parse()

	gpio.Setup()
	gpio.PinMode(FanControl, gpio.OUTPUT)
	gpio.PinMode(LedData, gpio.OUTPUT)
	gpio.PinMode(LedClock, gpio.OUTPUT)

	led := NewLed(0.1)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	gpio.DigitalWrite(FanControl, 1)

	time.Sleep(100 * time.Millisecond)

	var x = 1
	for {
		select {
		case <-signalChan:
			led.Clear()
			gpio.Cleanup()
			return
		default:
			temp, err := getCpuTemp()

			if err == nil {
				if temp < *offThreshold && x == 1 {
					x = 0
					gpio.DigitalWrite(FanControl, 0)
				}

				if temp >= *onThreshold && x == 0 {
					x = 1
					gpio.DigitalWrite(FanControl, 1)
				}

				updateLed(&led, temp)

				if *verbose {
					fmt.Printf("Current: %f, Target: %f, On: %t \n", temp, *offThreshold, x == 1)
				}
			} else {
				fmt.Println(err)
			}

			time.Sleep(time.Duration(*delay) * time.Second)
		}
	}
}

const minTemp = 35
const maxTemp = 80

func updateLed(led *Led, temp float64) {
	var hue float64

	lowTemp := *offThreshold
	highTemp := *onThreshold

	if temp < lowTemp && *extendedColors {
		// Between minimum temp and low temp, set LED to blue through to green
		temp = temp - minTemp
		temp = temp / (lowTemp - minTemp)
		temp = math.Max(0, temp)
		hue = (120.0 / 360.0) + ((1.0 - temp) * 120.0 / 360.0)
	} else if temp > highTemp && *extendedColors {
		// Between high temp and maximum temp, set LED to red through to magenta
		temp = temp - highTemp
		temp = temp / (maxTemp - highTemp)
		temp = math.Min(1, temp)
		hue = 1.0 - (temp * 60.0 / 360.0)
	} else {
		// In the normal low temp to high temp range, set LED to green through to red
		temp = temp - lowTemp
		temp = temp / (highTemp - lowTemp)
		temp = math.Max(0, math.Min(1, temp))
		hue = (1.0 - temp) * 120.0 / 360.0
	}

	color := colorful.Hsv(hue*360, 1.0, 1.0)
	led.SetPixel(int(color.R*255), int(color.G*255), int(color.B*255))
}

func getCpuTemp() (float64, error) {
	temperatures, _ := host.SensorsTemperatures()

	for _, s := range temperatures {
		if strings.HasPrefix(s.SensorKey, "cpu-thermal") || strings.HasPrefix(s.SensorKey, "cpu_thermal") {
			return s.Temperature, nil
		}
	}

	return -1, fmt.Errorf("unable to find CPU temperature")
}
