package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type PlatformBitmask struct {
	Enable1    bool
	Enable2    bool
	Enable3    bool
	Enable4    bool
	Enable5    bool
	Enable6    bool
	Enable7    bool
	EnableMore bool
}

// GetBits() returns the bitmask of the selected
// platform values
func (pb *PlatformBitmask) GetBits() int {
	total := 0

	if !pb.Enable1 {
		total += 1
	}
	if !pb.Enable2 {
		total += 2
	}
	if !pb.Enable3 {
		total += 4
	}
	if !pb.Enable4 {
		total += 8
	}
	if !pb.Enable5 {
		total += 16
	}
	if !pb.Enable6 {
		total += 32
	}
	if !pb.Enable7 {
		total += 64
	}
	if !pb.EnableMore {
		total += 128
	}

	return total
}

// CLass ID
type PlatformLayout struct {
	Platforms PlatformBitmask
	Lengths   PlatformBitmask
}

type AllowedPlatforms struct {
	Bitmask PlatformBitmask
}
type AllowedLengths struct {
	Bitmask PlatformBitmask
}

func (c *AllowedPlatforms) GetBytes() int {
	return 2
}

func (c *AllowedPlatforms) GetString() string {
	return fmt.Sprintf("0C %s ", bytes.GetByte(c.Bitmask.GetBits()))
}

func (c *AllowedLengths) GetBytes() int {
	return 2
}

func (c *AllowedLengths) GetString() string {
	return fmt.Sprintf("0D %s ", bytes.GetByte(c.Bitmask.GetBits()))
}
