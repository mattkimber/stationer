package sprites

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type Spritesets struct {
	NumSprites int
	NumSets    int
	ID         int
}

func (s *Spritesets) GetComment() string {
	return "Sprite sets"
}

func (s *Spritesets) GetLines() []string {
	len := 6
	if s.ID >= 255 {
		len = 8
	}
	result := fmt.Sprintf("* %d 01 %s 00 %s %s %s", len, bytes.GetByte(FEATURE_STATIONS), bytes.GetVariableByte(s.ID), bytes.GetByte(s.NumSets), bytes.GetByte(s.NumSprites))

	return []string{result}
}
