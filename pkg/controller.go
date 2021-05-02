package pkg

import (
	"fmt"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type FanshimController struct {
	on      bool
	fanshim Fanshim
	config  *Config
}

func NewController(config *Config, fanshim Fanshim) *FanshimController {
	fanshim.On()
	return &FanshimController{
		on:      true,
		fanshim: fanshim,
		config:  config,
	}
}

func (f *FanshimController) Update(temperature float64) {
	if temperature < f.config.OffThreshold && f.on {
		f.fanshim.Off()
		f.on = false
	}

	if temperature >= f.config.OnThreshold && !f.on {
		f.fanshim.On()
		f.on = true
	}

	color := f.calculateColor(temperature)
	f.fanshim.SetColor(color)

	if f.config.Verbose {
		fmt.Printf("Current: %f, Target: %f, Max: %f, On: %t , Color: %s\n",
			temperature,
			f.config.OffThreshold,
			f.config.OnThreshold,
			f.on,
			color.Hex(),
		)
	}
}

func (f *FanshimController) Cleanup() {
	f.fanshim.Cleanup()
}

func (f *FanshimController) calculateColor(temp float64) colorful.Color {
	var hue float64

	lowTemp := f.config.OffThreshold
	highTemp := f.config.OnThreshold

	if temp < lowTemp && f.config.ExtendedColors {
		// Between minimum temp and low temp, set LED to blue through to green
		temp = temp - MinTemp
		temp = temp / (lowTemp - MinTemp)
		temp = math.Max(0, temp)
		hue = (120.0 / 360.0) + ((1.0 - temp) * 120.0 / 360.0)
	} else if temp > highTemp && f.config.ExtendedColors {
		// Between high temp and maximum temp, set LED to red through to magenta
		temp = temp - highTemp
		temp = temp / (MaxTemp - highTemp)
		temp = math.Min(1, temp)
		hue = 1.0 - (temp * 60.0 / 360.0)
	} else {
		// In the normal low temp to high temp range, set LED to green through to red
		temp = temp - lowTemp
		temp = temp / (highTemp - lowTemp)
		temp = math.Max(0, math.Min(1, temp))
		hue = (1.0 - temp) * 120.0 / 360.0
	}

	return colorful.Hsv(hue*360, 1.0, f.config.Brightness/255.0)
}