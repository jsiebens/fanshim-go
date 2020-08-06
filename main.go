package main

import (
	"flag"
	"fmt"
	"github.com/ikester/gpio"
	"github.com/shirou/gopsutil/host"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const FanControl = 18

var offThreshold = flag.Float64("off-threshold", 55.0, "Temperature threshold in degrees C to disable fan")
var onThreshold = flag.Float64("on-threshold", 65.0, "Temperature threshold in degrees C to enable fan")
var delay = flag.Int("delay", 2, "Delay, in seconds, between temperature readings")
var verbose = flag.Bool("verbose", false, "Output temp and fan status messages")

func main() {
	flag.Parse()

	gpio.Setup()
	gpio.PinMode(FanControl, gpio.OUTPUT)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	gpio.DigitalWrite(FanControl, 1)

	time.Sleep(100 * time.Millisecond)

	var x = 1
	for {
		select {
		case <-signalChan:
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

func getCpuTemp() (float64, error) {
	temperatures, _ := host.SensorsTemperatures()

	for _, s := range temperatures {
		if strings.HasPrefix(s.SensorKey, "cpu-thermal") || strings.HasPrefix(s.SensorKey, "cpu_thermal") {
			return s.Temperature, nil
		}
	}

	return -1, fmt.Errorf("unable to find CPU temperature")
}
