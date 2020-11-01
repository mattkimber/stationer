package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type StationFenceCallback struct {
	SetID            int
	DefaultSpriteSet int
}

func (s *StationFenceCallback) GetComment() string {
	return "Callback for station fence choice"
}

func (s *StationFenceCallback) getAction(checkValue, ifTrueValue, ifFalseValue string, setIDOffset int) string {
	length := 18

	// 85 = access lowest word of station variable
	// 68 = station info of nearby tile
	// 1024 = get bit 10 (tile belongs to current station)
	return fmt.Sprintf("* %d 02 04 %s\n"+
		"    85 68 %s\n"+ // Check variable 68 for tile to -1/0
		"    00 %s\n"+ // mask bits
		"    01\n"+ // 1 non-default option
		"    %s %s %s\n"+
		"    %s\n", // 80 - set bit 15 to show this is a final return value
		length,
		bytes.GetByte(s.SetID+setIDOffset),
		checkValue,
		bytes.GetWord(65535),
		ifTrueValue,
		bytes.GetWord(65535), // If the tile is not a station the value of the lower bits is 0xFFFF
		bytes.GetWord(65535),
		ifFalseValue)
}

func (s *StationFenceCallback) getCallback(chainStart int) string {
	length := 17

	return fmt.Sprintf(
		"* %d 02 04 %s 85 0C 00 FF FF 01\n"+
			"    %s 00 14 00 14 00\n"+
			"    %s 00",
		length,
		bytes.GetByte(s.SetID), // The callback decider is given the SetID
		bytes.GetByte(chainStart),
		bytes.GetByte(s.DefaultSpriteSet),
	)
}

func (s *StationFenceCallback) GetLines() []string {
	return []string{
		s.getAction("10", "04 80", "00 80", 1),                                   // N: true, S: check
		s.getAction("10", "06 80", "02 80", 2),                                   // N: false, S: check
		s.getAction("F0", bytes.GetWord(s.SetID+2), bytes.GetWord(s.SetID+1), 3), // N: check, S: unknown
		s.getCallback(s.SetID + 3),
	}
}
