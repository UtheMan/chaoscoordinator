package disk

import "errors"

type Flags struct {
	Time           string
	FillPercentage int
}

func Fill(subID string, flags Flags) error {
	return errors.New("test")
}
