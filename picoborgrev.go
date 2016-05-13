// Package picoborgrev allows users to control pico borg reverse motor controllers over I2C using the gobot.io robotic framework.
// See: https://www.piborg.org/, https://gobot.io/
package picoborgrev

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
)

var _ gobot.Driver = (*Driver)(nil)

const (
	pwmMax              = 255  // Max pwm value
	i2cMaxLen           = 4    // Max len of I2C message
	i2cIDPicoborgRev    = 0x15 // i2c id of picoborg rev board
	commandSetLED       = 1    // Set the LED status
	commandGetLED       = 2    // Get the LED status
	commandSetAFWD      = 3    // Set motor 2 PWM rate in a forwards direction
	commandSetAREV      = 4    // Set motor 2 PWM rate in a reverse direction
	commandGetA         = 5    // Get motor 2 direction and PWM rate
	commandSetBFWD      = 6    // Set motor 1 PWM rate in a forwards direction
	commandSetBREV      = 7    // Set motor 1 PWM rate in a reverse direction
	commandGetB         = 8    // Get motor 1 direction and PWM rate
	commandAllOFF       = 9    // Switch everything off
	commandResetEPO     = 10   // Resets the EPO flag, use after EPO has been tripped and switch is now clear
	commandGetEPO       = 11   // Get the EPO latched flag
	commandSetEPOIgnore = 12   // Set the EPO ignored flag, allows the system to run without an EPO
	commandGetEPOIgnore = 13   // Get the EPO ignored flag
	commadGetDriveFault = 14   // Get the drive fault flag, indicates faults such as short-circuits and under voltage
	commandSetAllFWD    = 15   // Set all motors PWM rate in a forwards direction
	commandSetAllREV    = 16   // Set all motors PWM rate in a reverse direction
	commandSetFailsafe  = 17   // Set the failsafe flag, turns the motors off if communication is interrupted
	commandGetFailsafe  = 18   // Get the failsafe flag
	commandSetENCMode   = 19   // Set the board into encoder or speed mode
	commandGetENCMode   = 20   // Get the boards current mode, encoder or speed
	commandMoveAFWD     = 21   // Move motor 2 forward by n encoder ticks
	commandMoveAREV     = 22   // Move motor 2 reverse by n encoder ticks
	commandMoveBFWD     = 23   // Move motor 1 forward by n encoder ticks
	commandMoveBREV     = 24   // Move motor 1 reverse by n encoder ticks
	commandMoveAllFWD   = 25   // Move all motors forward by n encoder ticks
	commandMoveAllREV   = 26   // Move all motors reverse by n encoder ticks
	commandGetENCMoving = 27   // Get the status of encoders moving
	commandENCSpeed     = 28   // Set the maximum PWM rate in encoder mode
	commandGetENCSpeed  = 29   // Get the maximum PWM rate in encoder mode
	commandGetID        = 0x99 // Get the board identifier
	commandSetI2cAddr   = 0xAA // Set a new I2C address

	commandValueFWD = 1 // I2C value representing forward
	commandValueREV = 2 // I2C value representing reverse

	commandValueOn  = 1 // I2C value representing on
	commandValueOff = 0 // I2C value representing off

)

// Driver struct
type Driver struct {
	name       string
	connection i2c.I2c
	address    int
}

// NewDriver creates a new driver with specified name and i2c interface
func NewDriver(a i2c.I2c, name string, address int) *Driver {
	return &Driver{
		name:       name,
		connection: a,
		address:    address,
	}
}

// Name returns the name of the device
func (h *Driver) Name() string { return h.name }

// Connection returns the connection
func (h *Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initialized the picoborgrev
func (h *Driver) Start() (errs []error) {
	if err := h.connection.I2cStart(h.address); err != nil {
		return []error{err}
	}
	h.connection.I2cWrite(h.address, []byte{commandGetID})
	d, err := h.connection.I2cRead(h.address, i2cMaxLen)
	if err != nil {
		return []error{err}
	}

	if len(d) == i2cMaxLen {
		if d[1] != i2cIDPicoborgRev {
			err := fmt.Errorf("Found a device but it is not a PicoBorg Revers (ID %X instead of %X)", d[1], i2cIDPicoborgRev)
			return []error{err}
		}
	} else {
		err := fmt.Errorf("Device not found")
		return []error{err}
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *Driver) Halt() (errs []error) { return } //TODO: To implement

// SetMotorA generic set motor speed function
func (h *Driver) SetMotorA(power float32) error {
	var command byte
	var pwm int
	if power < 0 {
		command = commandSetAREV
		pwm = -int(pwmMax * power)
	} else {
		command = commandSetAFWD
		pwm = int(pwmMax * power)
	}

	err := h.connection.I2cWrite(h.address, []byte{command, byte(pwm)})
	if err != nil {
		return err
	}
	return nil
}

// SetMotorB generic set motor speed function
func (h *Driver) SetMotorB(power float32) error {
	var command byte
	var pwm int
	if power < 0 {
		command = commandSetBREV
		pwm = -int(pwmMax * power)
	} else {
		command = commandSetBFWD
		pwm = int(pwmMax * power)
	}

	err := h.connection.I2cWrite(h.address, []byte{command, byte(pwm)})
	if err != nil {
		return err
	}
	return nil
}
