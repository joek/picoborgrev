package picoborgrev_test

import (
	. "github.com/joek/picoborgrev"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Picoborgrev", func() {
	It("Creates a new Driver instance", func() {
		var d RevDriver
		d = NewDriver(newI2cTestAdaptor("adaptor"), "Test", 0x44)
		Ω(d).Should(BeAssignableToTypeOf(&Driver{}))
	})

	Describe("Driver", func() {
		var driver *Driver
		var adaptor *i2cTestAdaptor
		BeforeEach(func() {
			adaptor = newI2cTestAdaptor("adaptor")
			driver = NewDriver(adaptor, "test", 0x44)
		})

		It("should respond to getter", func() {
			Ω(driver.Name()).Should(Equal("test"))
			Ω(driver.Connection()).Should(BeAssignableToTypeOf(newI2cTestAdaptor("adaptor")))
		})

		It("should start", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}

			err := driver.Start()

			Ω(err).Should(BeNil())
			Ω(res).Should(Equal([]byte{0x99}))
		})

		It("should fail on wrong board", func() {
			adaptor.I2cReadImpl = func(i int, l int) ([]byte, error) {
				b := make([]byte, l, l)
				b[1] = 0x16
				return b, nil
			}

			err := driver.Start()

			Ω(err).Should(Not(BeNil()))
		})

		It("should fail on wrong response", func() {
			adaptor.I2cReadImpl = func(i int, l int) ([]byte, error) {
				b := make([]byte, l-1, l-1)
				b[1] = 0x15
				return b, nil
			}

			err := driver.Start()

			Ω(err).Should(Not(BeNil()))
		})

		It("should SetMotorA", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}

			driver.SetMotorA(0.25)

			Ω(res[0]).Should(Equal(byte(3)))
			Ω(res[1]).Should(Equal(byte(63)))

			driver.SetMotorA(-0.25)
			Ω(res[0]).Should(Equal(byte(4)))
			Ω(res[1]).Should(Equal(byte(63)))
		})
		It("should SetMotorB", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}

			driver.SetMotorB(0.25)

			Ω(res[0]).Should(Equal(byte(6)))
			Ω(res[1]).Should(Equal(byte(63)))

			driver.SetMotorB(-0.25)
			Ω(res[0]).Should(Equal(byte(7)))
			Ω(res[1]).Should(Equal(byte(63)))
		})
	})
})

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
func (t *i2cTestAdaptor) Name() string             { return t.name }
func (t *i2cTestAdaptor) Connect() (errs []error)  { return }
func (t *i2cTestAdaptor) Finalize() (errs []error) { return }

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
