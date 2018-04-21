// Package picoborgrev allows users to control pico borg reverse motor controllers over I2C using the gobot.io robotic framework.
// See: https://www.piborg.org/, https://gobot.io/
package picoborgrev

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

var _ gobot.Driver = (*Driver)(nil)

const (
	picoborgAddress     = 0x44 // Default address
	pwmMax              = 255  // Max pwm value
	i2cMaxLen           = 4    // Max len of I2C message
	i2cIDPicoborgRev    = 0x15 // i2c id of picoborg rev board
	commandSetLED       = 0x1  // Set the LED status
	commandGetLED       = 0x2  // Get the LED status
	commandSetAFWD      = 0x3  // Set motor 2 PWM rate in a forwards direction
	commandSetAREV      = 0x4  // Set motor 2 PWM rate in a reverse direction
	commandGetA         = 0x5  // Get motor 2 direction and PWM rate
	commandSetBFWD      = 0x6  // Set motor 1 PWM rate in a forwards direction
	commandSetBREV      = 0x7  // Set motor 1 PWM rate in a reverse direction
	commandGetB         = 0x8  // Get motor 1 direction and PWM rate
	commandAllOFF       = 0x9  // Switch everything off
	commandResetEPO     = 0x10 // Resets the EPO flag, use after EPO has been tripped and switch is now clear
	commandGetEPO       = 0x11 // Get the EPO latched flag
	commandSetEPOIgnore = 0x12 // Set the EPO ignored flag, allows the system to run without an EPO
	commandGetEPOIgnore = 0x13 // Get the EPO ignored flag
	commadGetDriveFault = 0x14 // Get the drive fault flag, indicates faults such as short-circuits and under voltage
	commandSetAllFWD    = 0x15 // Set all motors PWM rate in a forwards direction
	commandSetAllREV    = 0x16 // Set all motors PWM rate in a reverse direction
	commandSetFailsafe  = 0x17 // Set the failsafe flag, turns the motors off if communication is interrupted
	commandGetFailsafe  = 0x18 // Get the failsafe flag
	commandSetENCMode   = 0x19 // Set the board into encoder or speed mode
	commandGetENCMode   = 0x20 // Get the boards current mode, encoder or speed
	commandMoveAFWD     = 0x21 // Move motor 2 forward by n encoder ticks
	commandMoveAREV     = 0x22 // Move motor 2 reverse by n encoder ticks
	commandMoveBFWD     = 0x23 // Move motor 1 forward by n encoder ticks
	commandMoveBREV     = 0x24 // Move motor 1 reverse by n encoder ticks
	commandMoveAllFWD   = 0x25 // Move all motors forward by n encoder ticks
	commandMoveAllREV   = 0x26 // Move all motors reverse by n encoder ticks
	commandGetENCMoving = 0x27 // Get the status of encoders moving
	commandENCSpeed     = 0x28 // Set the maximum PWM rate in encoder mode
	commandGetENCSpeed  = 0x29 // Get the maximum PWM rate in encoder mode
	commandGetID        = 0x99 // Get the board identifier
	commandSetI2cAddr   = 0xAA // Set a new I2C address

	commandValueFWD = 0x1 // I2C value representing forward
	commandValueREV = 0x2 // I2C value representing reverse

	commandValueOn  = 0x1 // I2C value representing on
	commandValueOff = 0x0 // I2C value representing off

)

// RevDriver pico borg rev driver interace
type RevDriver interface {
	Name() string
	Connection() gobot.Connection
	Start() error
	Halt() error
	ResetEPO() error
	GetEPO() (bool, error)
	SetMotorA(float32) error
	SetMotorB(float32) error
	StopAllMotors() error
}

// Driver struct
type Driver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config
	lock sync.Mutex
}

// NewDriver creates a new driver with specified name and i2c interface
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
//
func NewDriver(a i2c.Connector, options ...func(i2c.Config)) *Driver {
	d := &Driver{
		name:      gobot.DefaultName("PicoBorg"),
		connector: a,
		Config:    i2c.NewConfig(),
		lock:      sync.Mutex{},
	}

	for _, option := range options {
		option(d)
	}
	return d
}

// Name returns the name of the device
func (h *Driver) Name() string { return h.name }

// SetName of the device
func (h *Driver) SetName(n string) { h.name = n }

// Connection returns the connection
func (h *Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initialized the picoborgrev
func (h *Driver) Start() (err error) {
	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())
	address := h.GetAddressOrDefault(picoborgAddress)
	h.connection, err = h.connector.GetConnection(address, bus)

	h.lock.Lock()
	defer h.lock.Unlock()
	if err != nil {
		return err
	}

	h.connection.Write([]byte{commandGetID})
	data := []byte{0, 0, 0, 0}

	read, err := h.connection.Read(data)
	if err != nil {
		return err
	}

	if read == i2cMaxLen {
		if data[1] != i2cIDPicoborgRev {
			err := fmt.Errorf("Found a device but it is not a PicoBorg Revers (ID %X instead of %X)", data[1], i2cIDPicoborgRev)
			return err
		}
	} else {
		err := fmt.Errorf("Device not found")
		return err
	}
	return
}

// Halt stops all motors
func (h *Driver) Halt() (errs error) {
	err := h.StopAllMotors()
	if err != nil {
		multierror.Append(errs, err)
	}
	return nil
}

// StopAllMotors will stop all motors
func (h *Driver) StopAllMotors() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	err := h.connection.WriteByte(byte(commandAllOFF))
	return err
}

// ResetEPO latch state, use to allow movement again after the EPO has been tripped
func (h *Driver) ResetEPO() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	err := h.connection.WriteByte(byte(commandResetEPO))
	return err
}

// GetEPO Reads the system EPO latch state.
func (h *Driver) GetEPO() (bool, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	err := h.connection.WriteByte(byte(commandGetEPO))
	if err != nil {
		return false, err
	}

	data, err := h.connection.ReadByte()
	if err != nil {
		return false, err
	}

	if int(data) == commandValueOff {
		return false, nil
	}
	return true, nil
}

// SetMotorA generic set motor speed function
func (h *Driver) SetMotorA(power float32) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	var command byte
	var pwm int
	if power < 0 {
		command = commandSetAREV
		pwm = -int(pwmMax * power)
	} else {
		command = commandSetAFWD
		pwm = int(pwmMax * power)
	}

	err := h.connection.WriteByteData(byte(command), byte(pwm))
	if err != nil {
		return err
	}
	return nil
}

// SetMotorB generic set motor speed function
func (h *Driver) SetMotorB(power float32) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	var command byte
	var pwm int
	if power < 0 {
		command = commandSetBREV
		pwm = -int(pwmMax * power)
	} else {
		command = commandSetBFWD
		pwm = int(pwmMax * power)
	}

	err := h.connection.WriteByteData(byte(command), byte(pwm))
	if err != nil {
		return err
	}
	return nil
}
