package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

// CLass ID
type GeneralFlag struct {
	HasCustomFoundations bool
	SpreadCargo bool
}

func (c *GeneralFlag) GetBytes() int {
	return 2
}

func (c *GeneralFlag) GetString() string {
	value := 0

	if c.SpreadCargo {
		value += 2
	}

	if c.HasCustomFoundations {
		value += 8
	}
	return fmt.Sprintf("13 %s ", bytes.GetByte(value))
}
