package picoborgrev_test

import (
	. "github.com/joek/picoborgrev"
	. "github.com/joek/picoborgrev/revtesthelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Picoborgrev", func() {
	It("Creates a new Driver instance", func() {
		var d RevDriver
		d = NewDriver(NewI2cTestAdaptor("adaptor"), "Test", 0x44)
		Ω(d).Should(BeAssignableToTypeOf(&Driver{}))
	})

	Describe("Driver", func() {
		var driver *Driver
		var adaptor *I2cTestAdaptor
		BeforeEach(func() {
			adaptor = NewI2cTestAdaptor("adaptor")
			driver = NewDriver(adaptor, "test", 0x44)
		})

		It("should respond to getter", func() {
			Ω(driver.Name()).Should(Equal("test"))
			Ω(driver.Connection()).Should(BeAssignableToTypeOf(NewI2cTestAdaptor("adaptor")))
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

		It("should get EPO state true", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}
			adaptor.I2cReadImpl = func(i int, l int) ([]byte, error) {
				b := make([]byte, l-1, l-1)
				b[1] = byte(0x1)
				return b, nil
			}

			b, _ := driver.GetEPO()

			Ω(b).Should(BeTrue())
			Ω(res[0]).Should(Equal(byte(0x11)))
		})

		It("should get EPO state false", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}
			adaptor.I2cReadImpl = func(i int, l int) ([]byte, error) {
				b := make([]byte, l-1, l-1)
				b[1] = byte(0x0)
				return b, nil
			}

			b, _ := driver.GetEPO()

			Ω(b).Should(BeFalse())
			Ω(res[0]).Should(Equal(byte(0x11)))
		})

		It("StopAllMotors hould Stop all motors at halt", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}

			driver.StopAllMotors()

			Ω(res[0]).Should(Equal(byte(9)))
		})

		It("Halt should Stop all motors at halt", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}

			driver.Halt()

			Ω(res[0]).Should(Equal(byte(9)))
		})

		It("should reset EPO", func() {
			var res []byte
			adaptor.I2cWriteImpl = func(i int, b []byte) error {
				res = b
				return nil
			}

			driver.ResetEPO()

			Ω(res[0]).Should(Equal(byte(0x10)))
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
