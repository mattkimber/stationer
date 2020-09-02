package nfo

import (
	"fmt"
	bytes2 "github.com/mattkimber/stationer/internal/bytes"
)

type StationSet struct {
	SetID int
	NumLittleSets int
	NumLotsSets int
	SpriteSets []int
}

func (s *StationSet) GetLines() []string {
	bytes := 5 + (len(s.SpriteSets) * 2)

	result := fmt.Sprintf("* %d 02 04 %s %s %s",
		bytes,
		bytes2.GetByte(s.SetID),
		bytes2.GetByte(s.NumLittleSets),
		bytes2.GetByte(s.NumLotsSets),
	)

	for _, set := range s.SpriteSets {
		result += fmt.Sprintf(" %s", bytes2.GetWord(set))
	}

	return []string { result }
}