package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type Spritesets struct {
	NumSprites int
	NumSets int
}

func (s *Spritesets) GetLines() []string {
	result := fmt.Sprintf("* 4 01 %s %s %s", bytes.GetByte(FEATURE_STATIONS), bytes.GetByte(s.NumSets), bytes.GetByte(s.NumSprites))

	return []string{result}
}