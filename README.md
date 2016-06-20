# Gobot driver for PicoBorg Reverse

Allows users to control [pico borg reverse motor controllers](https://www.piborg.org/picoborgrev/) over I2C using the [gobot.io](https://gobot.io) robotic framework.

[![GoDoc](https://godoc.org/github.com/joek/picoborgrev?status.svg)](http://godoc.org/github.com/joek/picoborgrev)
[![Travis](https://travis-ci.org/joek/picoborgrev.svg?branch=master)](https://travis-ci.org/joek/picoborgrev)

# Installation on Raspberry

The library was tested on Raspberry PI. It should run on other Platforms provided by [gobot.io](https://gobot.io/documentation/platforms/) too.

It is possible to build the project directly on the Raspberry. A documentation can be found in the gobot documentation. Nevertheless cross compiling on OS X is tested.

Using go1.6 just run:
```
$ env GOOS=linux GOARCH=arm GOARM=6 go build
```

And upload the binary to the raspberry
```
$ scp <binary> pi@192.168.1.xxx:/home/pi/
```

# How to use

Install missing libraries:
```
$ go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/platforms/i2c
$ go get github.com/joek/picoborgrev
```

The library is used like every gobot driver:
```
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
```

# Open Topics
The PicoBorg Reverse is offering a lot of functionality. Currently only the setMotor(A/B) functions are implemented. Based on the piborg python sample library the following functions are missing:

- SetMotors
- GetMotorA
- GetMotorB
- SetEpoIgnore
- GetEpoIgnore
- SetLED
- GetLED
- GetDriveFault
- SetEncoderMoveMode
- GetEncoderMoveMode
- SetEncoderSpeed
- GetEncoderSpeed
- EncoderMoveMotor1
- EncoderMoveMotor2
- EncoderMoveMotors
- IsEncoderMoving
- WaitWhileEncoderMoving

As a view functions are implemented later maybe somebody out of the community likes to contribute to the library to add some out of the list.


# Development

The driver is tested using the amazing [ginkgo](https://onsi.github.io/ginkgo/) and [gomega](https://onsi.github.io/gomega/) libraries. You can run the tests using ```go test```

# License
MIT License
