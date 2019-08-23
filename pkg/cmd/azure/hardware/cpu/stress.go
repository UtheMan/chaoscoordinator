package cpu

import "errors"

type Flags struct {
	Time string
}

func Stress(subID string, flags Flags) error {
	return errors.New("test")
}
