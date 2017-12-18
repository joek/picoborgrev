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
		d = NewDriver(NewI2cTestAdaptor(NewI2cFakeConnection()))
		Ω(d).Should(BeAssignableToTypeOf(&Driver{}))
	})

	Describe("Driver", func() {
		var connection *I2cFakeConnection
		var driver *Driver
		var adaptor *I2cTestAdaptor
		BeforeEach(func() {
			connection = NewI2cFakeConnection()
			adaptor = NewI2cTestAdaptor(connection)
			driver = NewDriver(adaptor)
		})

		It("should respond to getter", func() {
			connection.ReadImpl = func(p []byte) (n int, err error) {
				p[1] = 0x15
				return 4, nil
			}
			driver.Start()
			Ω(driver.Name()).Should(ContainSubstring("PicoBorg"))
			Ω(driver.Connection()).Should(BeAssignableToTypeOf(connection))
		})

		It("should start", func() {
			var data []byte
			connection.WriteImpl = func(b []byte) (int, error) {
				data = b
				return 1, nil
			}

			connection.ReadImpl = func(p []byte) (n int, err error) {
				p[1] = 0x15
				return 4, nil
			}
			err := driver.Start()

			Ω(err).Should(BeNil())
			Ω(len(data)).Should(Equal(1))
			Ω(data).Should(Equal([]byte{0x99}))
		})

		It("should fail on wrong board", func() {
			connection.ReadImpl = func(p []byte) (n int, err error) {
				p[1] = 0x16
				return 1, nil
			}

			err := driver.Start()

			Ω(err).ShouldNot(BeNil())
		})

		It("should fail on wrong response", func() {
			connection.ReadImpl = func(p []byte) (n int, err error) {
				p[1] = 0x15
				return 1, nil
			}

			err := driver.Start()

			Ω(err).ShouldNot(BeNil())
		})
		Describe("Start Driver", func() {
			BeforeEach(func() {
				connection.ReadImpl = func(p []byte) (n int, err error) {
					p[1] = 0x16
					return 1, nil
				}
				driver.Start()
			})

			It("should get EPO state true", func() {
				var data byte
				connection.WriteByteImpl = func(b byte) error {
					data = b
					return nil
				}
				connection.ReadByteImpl = func() (p byte, err error) {
					p = 0x1
					return p, nil
				}

				b, _ := driver.GetEPO()

				Ω(b).Should(BeTrue())
				Ω(data).Should(Equal(byte(0x11)))
			})

			It("should get EPO state false", func() {
				var data byte
				connection.WriteByteImpl = func(b byte) error {
					data = b
					return nil
				}
				connection.ReadByteImpl = func() (p byte, err error) {
					p = 0x0
					return p, nil
				}

				b, _ := driver.GetEPO()

				Ω(b).Should(BeFalse())
				Ω(data).Should(Equal(byte(0x11)))
			})

			It("StopAllMotors hould Stop all motors at halt", func() {
				var data byte
				connection.WriteByteImpl = func(b byte) error {
					data = b
					return nil
				}

				driver.StopAllMotors()

				Ω(data).Should(Equal(byte(9)))
			})

			It("Halt should Stop all motors at halt", func() {
				var data byte
				connection.WriteByteImpl = func(b byte) error {
					data = b
					return nil
				}

				driver.Halt()

				Ω(data).Should(Equal(byte(9)))
			})

			It("should reset EPO", func() {
				var data byte
				connection.WriteByteImpl = func(b byte) error {
					data = b
					return nil
				}

				driver.ResetEPO()

				Ω(data).Should(Equal(byte(0x10)))
			})

			It("should SetMotorA", func() {
				var data uint8
				var c uint8
				connection.WriteByteDataImpl = func(reg uint8, val uint8) error {
					data = val
					c = reg
					return nil
				}

				driver.SetMotorA(0.25)

				Ω(c).Should(Equal(byte(3)))
				Ω(data).Should(Equal(byte(63)))

				driver.SetMotorA(-0.25)
				Ω(c).Should(Equal(byte(4)))
				Ω(data).Should(Equal(byte(63)))
			})
			It("should SetMotorB", func() {
				var data uint8
				var c uint8
				connection.WriteByteDataImpl = func(reg uint8, val uint8) error {
					data = val
					c = reg
					return nil
				}

				driver.SetMotorB(0.25)

				Ω(c).Should(Equal(byte(6)))
				Ω(data).Should(Equal(byte(63)))

				driver.SetMotorB(-0.25)
				Ω(c).Should(Equal(byte(7)))
				Ω(data).Should(Equal(byte(63)))
			})
		})
	})
})
