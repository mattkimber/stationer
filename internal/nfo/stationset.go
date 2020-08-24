package nfo

import "fmt"

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
		GetByte(s.SetID),
		GetByte(s.NumLittleSets),
		GetByte(s.NumLotsSets),
	)

	for _, set := range s.SpriteSets {
		result += fmt.Sprintf(" %s", GetShort(set))
	}

	return []string { result }
}