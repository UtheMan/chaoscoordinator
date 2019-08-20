package disk

import "errors"

type Flags struct {
	Time string
}

func Fill(subID string, flags Flags) error {
	return errors.New("test")
}
