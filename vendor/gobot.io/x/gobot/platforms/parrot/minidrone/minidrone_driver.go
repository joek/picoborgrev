package minidrone

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// Driver is gobot software device to the keyboard
type Driver struct {
	name       string
	connection gobot.Connection
	stepsfa0a  uint16
	stepsfa0b  uint16
	flying     bool
	Pcmd       Pcmd
	gobot.Eventer
}

const (
	// service IDs
	DroneCommandService      = "9a66fa000800919111e4012d1540cb8e"
	DroneNotificationService = "9a66fb000800919111e4012d1540cb8e"

	// characteristic IDs
	PcmdCharacteristic         = "9a66fa0a0800919111e4012d1540cb8e"
	CommandCharacteristic      = "9a66fa0b0800919111e4012d1540cb8e"
	FlightStatusCharacteristic = "9a66fb0e0800919111e4012d1540cb8e"
	BatteryCharacteristic      = "9a66fb0f0800919111e4012d1540cb8e"

	// Battery event
	Battery = "battery"

	// flight status event
	Status = "status"

	// flying event
	Flying = "flying"

	// landed event
	Landed = "landed"
)

type Pcmd struct {
	Flag  int
	Roll  int
	Pitch int
	Yaw   int
	Gaz   int
	Psi   float32
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	} else if val < 0 {
		return 0
	}

	return val
}

// NewDriver creates a Parrot Minidrone Driver
func NewDriver(a *ble.ClientAdaptor) *Driver {
	n := &Driver{
		name:       "Minidrone",
		connection: a,
		Pcmd: Pcmd{
			Flag:  0,
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
			Gaz:   0,
			Psi:   0,
		},
		Eventer: gobot.NewEventer(),
	}

	n.AddEvent(Battery)
	n.AddEvent(Status)
	n.AddEvent(Flying)
	n.AddEvent(Landed)

	return n
}

func (b *Driver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *Driver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *Driver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *Driver) adaptor() *ble.ClientAdaptor {
	return b.Connection().(*ble.ClientAdaptor)
}

// Start tells driver to get ready to do work
func (b *Driver) Start() (err error) {
	b.Init()
	b.FlatTrim()
	b.StartPcmd()
	b.FlatTrim()

	return
}

// Halt stops minidrone driver (void)
func (b *Driver) Halt() (err error) {
	b.Land()

	time.Sleep(500 * time.Millisecond)
	return
}

func (b *Driver) Init() (err error) {
	b.GenerateAllStates()

	// subscribe to battery notifications
	b.adaptor().Subscribe(DroneNotificationService, BatteryCharacteristic, func(data []byte, e error) {
		b.Publish(b.Event(Battery), data[len(data)-1])
	})

	// subscribe to flying status notifications
	b.adaptor().Subscribe(DroneNotificationService, FlightStatusCharacteristic, func(data []byte, e error) {
		if len(data) < 7 || data[2] != 2 {
			fmt.Println(data)
			return
		}
		b.Publish(b.Event(Status), data[6])
		if (data[6] == 1 || data[6] == 2) && !b.flying {
			b.flying = true
			b.Publish(b.Event(Flying), true)
		} else if (data[6] == 0) && b.flying {
			b.flying = false
			b.Publish(b.Event(Landed), true)
		}
	})

	return
}

func (b *Driver) GenerateAllStates() (err error) {
	b.stepsfa0b++
	buf := []byte{0x04, byte(b.stepsfa0b), 0x00, 0x04, 0x01, 0x00, 0x32, 0x30, 0x31, 0x34, 0x2D, 0x31, 0x30, 0x2D, 0x32, 0x38, 0x00}
	err = b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, buf)
	if err != nil {
		fmt.Println("GenerateAllStates error:", err)
		return err
	}

	return
}

func (b *Driver) TakeOff() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x00, 0x01, 0x00}
	err = b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, buf)
	if err != nil {
		fmt.Println("takeoff error:", err)
		return err
	}

	return
}

func (b *Driver) Land() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b), 0x02, 0x00, 0x03, 0x00}
	err = b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, buf)

	return err
}

func (b *Driver) FlatTrim() (err error) {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x00, 0x00, 0x00}
	err = b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, buf)

	return err
}

func (b *Driver) StartPcmd() {
	go func() {
		// wait a little bit so that there is enough time to get some ACKs
		time.Sleep(500 * time.Millisecond)
		for {
			err := b.adaptor().WriteCharacteristic(DroneCommandService, PcmdCharacteristic, b.generatePcmd().Bytes())
			if err != nil {
				fmt.Println("pcmd write error:", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (b *Driver) Up(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Gaz = validatePitch(val)
	return nil
}

func (b *Driver) Down(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Gaz = validatePitch(val) * -1
	return nil
}

func (b *Driver) Forward(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Pitch = validatePitch(val)
	return nil
}

func (b *Driver) Backward(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Pitch = validatePitch(val) * -1
	return nil
}

func (b *Driver) Right(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Roll = validatePitch(val)
	return nil
}

func (b *Driver) Left(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Roll = validatePitch(val) * -1
	return nil
}

func (b *Driver) Clockwise(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Yaw = validatePitch(val)
	return nil
}

func (b *Driver) CounterClockwise(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Yaw = validatePitch(val) * -1
	return nil
}

func (b *Driver) Stop() error {
	b.Pcmd = Pcmd{
		Flag:  0,
		Roll:  0,
		Pitch: 0,
		Yaw:   0,
		Gaz:   0,
		Psi:   0,
	}

	return nil
}

// StartRecording not supported
func (b *Driver) StartRecording() error {
	return nil
}

// StopRecording not supported
func (b *Driver) StopRecording() error {
	return nil
}

// HullProtection not supported
func (b *Driver) HullProtection(protect bool) error {
	return nil
}

// Outdoor not supported
func (b *Driver) Outdoor(outdoor bool) error {
	return nil
}

func (b *Driver) FrontFlip() (err error) {
	return b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, b.generateAnimation(0).Bytes())
}

func (b *Driver) BackFlip() (err error) {
	return b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, b.generateAnimation(1).Bytes())
}

func (b *Driver) RightFlip() (err error) {
	return b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, b.generateAnimation(2).Bytes())
}

func (b *Driver) LeftFlip() (err error) {
	return b.adaptor().WriteCharacteristic(DroneCommandService, CommandCharacteristic, b.generateAnimation(3).Bytes())
}

func (b *Driver) generateAnimation(direction int8) *bytes.Buffer {
	b.stepsfa0b++
	buf := []byte{0x02, byte(b.stepsfa0b) & 0xff, 0x02, 0x04, 0x00, 0x00, byte(direction), 0x00, 0x00, 0x00}
	return bytes.NewBuffer(buf)
}

func (b *Driver) generatePcmd() *bytes.Buffer {
	b.stepsfa0a++

	cmd := &bytes.Buffer{}
	binary.Write(cmd, binary.LittleEndian, int8(2))
	binary.Write(cmd, binary.LittleEndian, int8(b.stepsfa0a))
	binary.Write(cmd, binary.LittleEndian, int8(2))
	binary.Write(cmd, binary.LittleEndian, int8(0))
	binary.Write(cmd, binary.LittleEndian, int8(2))
	binary.Write(cmd, binary.LittleEndian, int8(0))
	binary.Write(cmd, binary.LittleEndian, int8(b.Pcmd.Flag))
	binary.Write(cmd, binary.LittleEndian, int8(b.Pcmd.Roll))
	binary.Write(cmd, binary.LittleEndian, int8(b.Pcmd.Pitch))
	binary.Write(cmd, binary.LittleEndian, int8(b.Pcmd.Yaw))
	binary.Write(cmd, binary.LittleEndian, int8(b.Pcmd.Gaz))
	binary.Write(cmd, binary.LittleEndian, float32(b.Pcmd.Psi))
	binary.Write(cmd, binary.LittleEndian, int16(0))
	binary.Write(cmd, binary.LittleEndian, int16(0))

	return cmd
}
