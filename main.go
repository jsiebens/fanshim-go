package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/jasonlvhit/gocron"
	"github.com/jsiebens/fanshim-go/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/ztrue/shutdown"
)

func main() {
	config := pkg.LoadConfig()

	fanshim := pkg.NewGPIOFanshim()
	controller := pkg.NewController(config, fanshim)

	shutdown.Add(func() {
		gocron.Clear()
		controller.Cleanup()
	})

	gocron.Every(config.Delay).Seconds().Do(tick, controller)
	gocron.Start()

	go func() {
		fmt.Println("Starting metrics http service ...")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Starting shutdown listener ...")
	shutdown.Listen(os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
}

func tick(controller *pkg.FanshimController) {
	temp, err := getCpuTemp()
	if err != nil {
		fmt.Println(err)
		return
	}
	percent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	controller.Update(temp, percent[0])
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
