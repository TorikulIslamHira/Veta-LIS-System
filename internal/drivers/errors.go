package drivers

import "fmt"

var (
	ErrUnsupportedMachineType = fmt.Errorf("unsupported machine type")
	ErrConnectionFailed       = fmt.Errorf("connection failed")
	ErrAlreadyConnected       = fmt.Errorf("already connected")
	ErrNotConnected           = fmt.Errorf("not connected")
	ErrInvalidData            = fmt.Errorf("invalid data format")
	ErrParseError             = fmt.Errorf("data parse error")
	ErrCommandTimeout         = fmt.Errorf("command timeout")
)
