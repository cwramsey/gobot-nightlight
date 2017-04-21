package main

import (
	"time"

	"fmt"

	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func info(msg ...interface{}) {
	fmt.Printf("[INFO] %v - %v\n", time.Now(), msg)
}

func main() {
	// Setup the initial state of the system
	state := "off"
	waitTime := 15 * time.Minute
	var lastStartTime time.Time

	// setup all the pins
	switchPin := os.Getenv("SWITCHPIN")
	lightPin := os.Getenv("LIGHTPIN")
	buttonPin := os.Getenv("BUTTONPIN")

	if switchPin == "" || lightPin == "" || buttonPin == "" {
		fmt.Println("The following environment variables are required. SWITCHPIN, LIGHTPIN, BUTTONPIN")
	}

	r := raspi.NewAdaptor()
	powerSwitch := gpio.NewLedDriver(r, switchPin)
	light := gpio.NewLedDriver(r, "8")
	button := gpio.NewButtonDriver(r, "15")

	// turn on the power switch's LED
	powerSwitch.On()
	powerSwitch.Brightness(255)

	// get started working!
	work := func() {

		// when the button is pushed, toggle the light on or off.
		button.Once("push", func(data interface{}) {
			info("button pushed. state:", state, "data:", data)
			switch state {
			case "off":
				state = "on"
				light.On()
				lastStartTime = time.Now()

			case "on":
				state = "off"
				light.Off()
			}

		})

		gobot.Every(50*time.Millisecond, func() {
			if time.Now().Sub(lastStartTime) > waitTime {
				state = "off"
				light.Off()
			}
		})
	}

	robot := gobot.NewRobot("nightlight",
		[]gobot.Connection{r},
		[]gobot.Device{powerSwitch},
		work,
	)

	robot.Start()
}
