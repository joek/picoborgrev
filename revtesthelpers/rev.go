package revtesthelpers

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

type I2cTestAdaptor struct {
	GetConnectionImpl func(address int, bus int) (device i2c.Connection, err error)
}

func (t *I2cTestAdaptor) GetConnection(address int, bus int) (device i2c.Connection, err error) {
	return t.GetConnectionImpl(address, bus)
}
func (t *I2cTestAdaptor) GetDefaultBus() int { return 1 }

func NewI2cTestAdaptor(conn *I2cFakeConnection) *I2cTestAdaptor {
	return &I2cTestAdaptor{
		GetConnectionImpl: func(address int, bus int) (device i2c.Connection, err error) {
			device = conn
			return
		},
	}
}

type I2cFakeConnection struct {
	CloseImpl          func() error
	ReadImpl           func(p []byte) (n int, err error)
	WriteImpl          func(p []byte) (n int, err error)
	ReadByteImpl       func() (val byte, err error)
	ReadByteDataImpl   func(reg uint8) (val uint8, err error)
	ReadWordDataImpl   func(reg uint8) (val uint16, err error)
	WriteByteImpl      func(val byte) (err error)
	WriteByteDataImpl  func(reg uint8, val uint8) (err error)
	WriteWordDataImpl  func(reg uint8, val uint16) (err error)
	WriteBlockDataImpl func(reg uint8, b []byte) (err error)
	ConnectImpl        func() error
	FinalizeImpl       func() error
	NameImpl           func() string
	SetNameImpl        func(string)
}

func (c *I2cFakeConnection) Close() error { return c.CloseImpl() }
func (c *I2cFakeConnection) Read(p []byte) (n int, err error) {
	return c.ReadImpl(p)
}
func (c *I2cFakeConnection) Write(p []byte) (n int, err error) { return c.WriteImpl(p) }
func (c *I2cFakeConnection) ReadByte() (val byte, err error)   { return c.ReadByteImpl() }
func (c *I2cFakeConnection) ReadByteData(reg uint8) (val uint8, err error) {
	return c.ReadByteDataImpl(reg)
}
func (c *I2cFakeConnection) ReadWordData(reg uint8) (val uint16, err error) {
	return c.ReadWordDataImpl(reg)
}
func (c *I2cFakeConnection) WriteByte(val byte) (err error) { return c.WriteByteImpl(val) }
func (c *I2cFakeConnection) WriteByteData(reg uint8, val uint8) (err error) {
	return c.WriteByteDataImpl(reg, val)
}
func (c *I2cFakeConnection) WriteWordData(reg uint8, val uint16) (err error) {
	return c.WriteWordDataImpl(reg, val)
}
func (c *I2cFakeConnection) WriteBlockData(reg uint8, b []byte) (err error) {
	return c.WriteBlockDataImpl(reg, b)
}
func (c *I2cFakeConnection) Connect() error {
	return c.ConnectImpl()
}
func (c *I2cFakeConnection) Name() string {
	return c.NameImpl()
}
func (c *I2cFakeConnection) SetName(n string) {
	c.SetNameImpl(n)
}
func (c *I2cFakeConnection) Finalize() error {
	return c.FinalizeImpl()
}

func NewI2cFakeConnection() *I2cFakeConnection {
	c := &I2cFakeConnection{
		WriteImpl: func(b []byte) (int, error) {
			return 1, nil
		},
	}
	return c
}

type FakeRevDriver struct {
	name              string
	connection        gobot.Connection
	SetMotorAImpl     func(float32) error
	SetMotorBImpl     func(float32) error
	StartImpl         func() error
	StopAllMotorsImpl func() error
	HaltImpl          func() error
	ResetEPOImpl      func() error
	GetEPOImpl        func() (bool, error)
}

func NewFakeRevDriver() *FakeRevDriver {
	return &FakeRevDriver{
		name:       "FakeRevDriver",
		connection: newI2cTestAdaptor("I2CTest"),
		SetMotorAImpl: func(power float32) error {
			return nil
		},
		SetMotorBImpl: func(power float32) error {
			return nil
		},
		StartImpl: func() error {
			return nil
		},
		StopAllMotorsImpl: func() error {
			return nil
		},
		HaltImpl: func() error {
			return nil
		},
		GetEPOImpl: func() (bool, error) {
			return true, nil
		},
		ResetEPOImpl: func() error {
			return nil
		},
	}
}

func (b *FakeRevDriver) SetName(n string) {
	b.name = n
}

func (b *FakeRevDriver) SetMotorA(power float32) error {
	return b.SetMotorAImpl(power)
}

func (b *FakeRevDriver) SetMotorB(power float32) error {
	return b.SetMotorBImpl(power)
}

func (b *FakeRevDriver) Start() error {
	return b.StartImpl()
}

func (b *FakeRevDriver) Halt() error {
	return b.HaltImpl()
}

func (b *FakeRevDriver) Name() string {
	return b.name
}

func (b *FakeRevDriver) Connection() gobot.Connection {
	return b.connection.(gobot.Connection)
}

func (b *FakeRevDriver) ResetEPO() error {
	return b.ResetEPOImpl()
}
func (b *FakeRevDriver) GetEPO() (bool, error) {
	return b.GetEPOImpl()
}

func (b *FakeRevDriver) StopAllMotors() error {
	return b.StopAllMotorsImpl()
}

type i2cTestAdaptor struct {
	name         string
	I2cReadImpl  func(i int, l int) ([]byte, error)
	I2cWriteImpl func(int, []byte) error
	I2cStartImpl func() error
}

func (t *i2cTestAdaptor) I2cStart(int) (err error) {
	return t.I2cStartImpl()
}
func (t *i2cTestAdaptor) I2cRead(i int, l int) (data []byte, err error) {
	return t.I2cReadImpl(i, l)
}
func (t *i2cTestAdaptor) I2cWrite(i int, b []byte) (err error) {
	return t.I2cWriteImpl(i, b)
}
func (t *i2cTestAdaptor) Name() string           { return t.name }
func (t *i2cTestAdaptor) SetName(n string)       { t.name = n }
func (t *i2cTestAdaptor) Connect() (errs error)  { return }
func (t *i2cTestAdaptor) Finalize() (errs error) { return }

func newI2cTestAdaptor(name string) *i2cTestAdaptor {
	return &i2cTestAdaptor{
		name: name,
		I2cReadImpl: func(i int, l int) ([]byte, error) {
			b := make([]byte, l, l)
			b[1] = 0x15
			return b, nil
		},
		I2cWriteImpl: func(i int, b []byte) error {
			return nil
		},
		I2cStartImpl: func() error {
			return nil
		},
	}
}
