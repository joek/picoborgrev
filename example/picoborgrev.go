package main

import (
	"time"

	"github.com/joek/picoborgrev"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	rev := picoborgrev.NewDriver(r, "rev", 10)

	work := func() {
		rev.SetMotorA(0.3)
		time.Sleep(10 * time.Second)
		rev.SetMotorA(-0.3)
		time.Sleep(10 * time.Second)
		rev.SetMotorA(0)
	}

	robot := gobot.NewRobot("beerbot",
		[]gobot.Connection{r},
		[]gobot.Device{rev},
		work,
	)

	robot.Start()
}
