package pkg

import (
	"fmt"

	"github.com/ikester/gpio"
	"github.com/lucasb-eyer/go-colorful"
)

const FanControl = 18
const LedData int = 15
const LedClock int = 14

type Fanshim interface {
	On()
	Off()
	SetColor(color colorful.Color)
	Cleanup()
}

func NewGPIOFanshim() Fanshim {
	g := &GPIOFanshim{}
	g.init()
	return g
}

type GPIOFanshim struct {
	led Led
}

func (f *GPIOFanshim) init() {
	gpio.Setup()
	gpio.PinMode(FanControl, gpio.OUTPUT)
	gpio.PinMode(LedData, gpio.OUTPUT)
	gpio.PinMode(LedClock, gpio.OUTPUT)
	f.led = newLed(0.1)
}

func (f *GPIOFanshim) On() {
	gpio.DigitalWrite(FanControl, 1)
}

func (f *GPIOFanshim) Off() {
	gpio.DigitalWrite(FanControl, 0)
}

func (f *GPIOFanshim) SetColor(color colorful.Color) {
	f.led.SetPixel(int(color.R*255), int(color.G*255), int(color.B*255))
}

func (f *GPIOFanshim) Cleanup() {
	f.led.Clear()
	gpio.Cleanup()
}

func NewFakeFanshim() Fanshim {
	return &FakeFanshim{}
}

type FakeFanshim struct {
}

func (f *FakeFanshim) On() {
	fmt.Println("Turning fan on!")
}

func (f *FakeFanshim) Off() {
	fmt.Println("Turning fan off!")
}

func (f *FakeFanshim) SetColor(color colorful.Color) {
	fmt.Printf("Updating color to %s\n", color.Hex())
}

func (f *FakeFanshim) Cleanup() {
	fmt.Println("Cleaning up")
}
