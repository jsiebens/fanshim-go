package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jsiebens/fanshim-go/pkg"
	"github.com/shirou/gopsutil/host"
)

func main() {
	config := pkg.LoadConfig()

	fanshim := pkg.NewGPIOFanshim()
	controller := pkg.NewController(config, fanshim)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	time.Sleep(100 * time.Millisecond)

	for {
		select {
		case <-signalChan:
			controller.Cleanup()
			return
		default:
			temp, err := getCpuTemp()

			if err == nil {
				controller.Update(temp)
			} else {
				fmt.Println(err)
			}

			time.Sleep(time.Duration(config.Delay) * time.Second)
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
