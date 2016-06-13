package revtesthelpers

type I2cTestAdaptor struct {
	name         string
	I2cReadImpl  func(i int, l int) ([]byte, error)
	I2cWriteImpl func(int, []byte) error
	I2cStartImpl func() error
}

func (t *I2cTestAdaptor) I2cStart(int) (err error) {
	return t.I2cStartImpl()
}
func (t *I2cTestAdaptor) I2cRead(i int, l int) (data []byte, err error) {
	return t.I2cReadImpl(i, l)
}
func (t *I2cTestAdaptor) I2cWrite(i int, b []byte) (err error) {
	return t.I2cWriteImpl(i, b)
}
func (t *I2cTestAdaptor) Name() string             { return t.name }
func (t *I2cTestAdaptor) Connect() (errs []error)  { return }
func (t *I2cTestAdaptor) Finalize() (errs []error) { return }

func NewI2cTestAdaptor(name string) *I2cTestAdaptor {
	return &I2cTestAdaptor{
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
