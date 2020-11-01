package nfo

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
	result := fmt.Sprintf("* 6 01 %s 00 %s %s %s", bytes.GetByte(FEATURE_STATIONS), bytes.GetByte(s.ID), bytes.GetByte(s.NumSets), bytes.GetByte(s.NumSprites))

	return []string{result}
}
