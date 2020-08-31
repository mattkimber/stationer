package nfo

import "fmt"

type Spritesets struct {
	NumSprites int
	NumSets int
}

func (s *Spritesets) GetLines() []string {
	result := fmt.Sprintf("* 4 01 %s %s %s", GetByte(FEATURE_STATIONS), GetByte(s.NumSets), GetByte(s.NumSprites))

	return []string{result}
}