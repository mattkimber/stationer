package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

// CLass ID
type CallbackFlag struct {
	SpriteLayout bool
}

func (c *CallbackFlag) GetBytes() int {
	return 2
}

func (c *CallbackFlag) GetString() string {
	value := 0
	if c.SpriteLayout {
		value += 2
	}
	return fmt.Sprintf("0B %s ", bytes.GetByte(value))
}
