package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type BlockedPillarInformation struct {
	IsBlocked bool
	Layouts   int
}

func (b *BlockedPillarInformation) GetBytes() int {
	return 2 + b.Layouts
}

func (b *BlockedPillarInformation) GetString() string {
	flags := ""
	for i := 0; i < b.Layouts; i++ {
		if b.IsBlocked {
			flags += fmt.Sprintf("01 ")
		} else {
			flags += fmt.Sprintf("00 ")
		}
	}

	return fmt.Sprintf("21 %s %s // bridge pillar flags", bytes.GetByte(b.Layouts), flags)
}
