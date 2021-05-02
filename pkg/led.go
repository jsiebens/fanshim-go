package pkg

import (
	"log"

	"github.com/ikester/gpio"
)

const redIndex int = 0
const greenIndex int = 1
const blueIndex int = 2
const brightnessIndex int = 3

// default raw brightness.  Not to be used user-side
const defaultBrightnessInt int = 10

//upper and lower bounds for user specified brightness
const minBrightness float64 = 0.0
const maxBrightness float64 = 1.0

func newLed(brightness ...float64) Led {
	brightnessInt := defaultBrightnessInt
	if len(brightness) > 0 {
		brightnessInt = convertBrightnessToInt(brightness[0])
	}
	return Led{
		pixel: initPixel(brightnessInt),
	}
}

type Led struct {
	pixel [4]int
}

// pulse sends a pulse through the LedData/CLK pins
func pulse(pulses int) {
	gpio.DigitalWrite(LedData, 0)
	for i := 0; i < pulses; i++ {
		gpio.DigitalWrite(LedClock, 1)
		gpio.DigitalWrite(LedClock, 0)
	}
}

// eof end of file or signal, from Python library
func eof() {
	pulse(36)
}

// sof start of file (name from Python library)
func sof() {
	pulse(32)
}

func writeByte(val int) {
	for i := 0; i < 8; i++ {
		// 0b10000000 = 128
		gpio.DigitalWrite(LedData, val&128)
		gpio.DigitalWrite(LedClock, 1)
		val = val << 1
		gpio.DigitalWrite(LedClock, 0)
	}
}

func convertBrightnessToInt(brightness float64) int {
	if !inRangeFloat(minBrightness, brightness, maxBrightness) {
		log.Fatalf("Supplied brightness was %#v - value should be between: %#v and %#v", brightness, minBrightness, maxBrightness)
	}
	return int(brightness * 31.0)
}

func inRangeFloat(minVal float64, testVal float64, maxVal float64) bool {
	return (testVal >= minVal) && (testVal <= maxVal)
}

func (bl *Led) Clear() {
	bl.SetPixel(0, 0, 0)
	bl.show()
}

func (bl *Led) show() {
	sof()
	brightness := bl.pixel[brightnessIndex]
	r := bl.pixel[redIndex]
	g := bl.pixel[greenIndex]
	b := bl.pixel[blueIndex]

	// 0b11100000 (224)
	bitwise := 224
	writeByte(bitwise | brightness)
	writeByte(b)
	writeByte(g)
	writeByte(r)
	eof()
}

func (bl *Led) SetPixel(r int, g int, b int) {
	bl.pixel[redIndex] = r
	bl.pixel[greenIndex] = g
	bl.pixel[blueIndex] = b
	bl.show()
}

func (bl *Led) SetBrightness(brightness float64) *Led {
	brightnessInt := convertBrightnessToInt(brightness)
	bl.pixel[brightnessIndex] = brightnessInt

	bl.show()

	return bl
}

func initPixel(brightness int) [4]int {
	var pixels [4]int
	pixels[redIndex] = 0
	pixels[greenIndex] = 0
	pixels[blueIndex] = 0
	pixels[brightnessIndex] = brightness
	return pixels
}
