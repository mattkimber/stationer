package properties

import "fmt"

type PreventTrainEntryFlag struct {
	ID string
}

func (c *PreventTrainEntryFlag) GetBytes() int {
	return 2
}

func (c *PreventTrainEntryFlag) GetString() string {
	return fmt.Sprintf("15 FF ")
}
