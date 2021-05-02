package pkg

import (
	"flag"
	"os"
	"strconv"
)

const MinTemp = 35
const MaxTemp = 80

type Config struct {
	OffThreshold   float64
	OnThreshold    float64
	ExtendedColors bool
	Delay          int
	Brightness     float64
	Verbose        bool
}

func LoadConfig() *Config {
	offThreshold := flag.Int("off-threshold", 55.0, "Temperature threshold in degrees C to disable fan")
	onThreshold := flag.Int("on-threshold", 65.0, "Temperature threshold in degrees C to enable fan")
	extendedColors := flag.Bool("extended-colors", false, "Extend LED colors for outside of normal low to high range")
	delay := flag.Int("delay", 2, "Delay, in seconds, between temperature readings")
	brightness := flag.Int("brightness", 255, "LED brightness, from 0 to 255")
	verbose := flag.Bool("verbose", false, "Output temp and fan status messages")

	flag.Parse()

	return &Config{
		OffThreshold:   float64(intValue("OFF_THRESHOLD", *offThreshold)),
		OnThreshold:    float64(intValue("ON_THRESHOLD", *onThreshold)),
		ExtendedColors: boolValue("EXTENDED_COLORS", *extendedColors),
		Delay:          intValue("DELAY", *delay),
		Brightness:     float64(intValue("BRIGHTNESS", *brightness)),
		Verbose:        boolValue("VERBOSE", *verbose),
	}
}

func boolValue(name string, fallback bool) bool {
	val := os.Getenv(name)
	if len(val) > 0 {
		return val == "true"
	}
	return fallback
}

func intValue(name string, fallback int) int {
	val := os.Getenv(name)
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return parsedVal
		}
	}
	return fallback
}
