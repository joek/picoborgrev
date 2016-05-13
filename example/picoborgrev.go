package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/joek/picoborgrev"
)

func main() {
	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
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

	gbot.AddRobot(robot)

	gbot.Start()
}
