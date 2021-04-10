package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

const (
	STATION_FENCE_OFFSET = 3
	FIRST_CAR_PARK_ID = 212
	FIRST_REVERSE_CAR_PARK_ID = 224
	NUM_CAR_PARKS = 4
)

type StationFenceCallback struct {
	SetID            int
	BaseLayoutOffset int
	DefaultSpriteSet int
	YearCallbackID   int
	HasDecider       bool
}

func (s *StationFenceCallback) GetComment() string {
	return "Callback for station fence choice"
}

func (s *StationFenceCallback) getAction(checkValue, ifTrueValue, ifFalseValue string, setIDOffset int) string {
	length := 30

	// 85 = access lowest word of station variable
	// 68 = station info of nearby tile
	// 1024 = get bit 10 (tile belongs to current station)
	return fmt.Sprintf("* %d 02 04 %s\n"+
		"    85 68 %s\n"+ // Check variable 68 for tile to -1/0
		"    00 %s\n"+ // mask bits
		"    03\n"+ // 3 non-default options
		"    %s %s %s\n"+ // Everything is FF
		"    %s %s %s\n" + // Line for car parks to still get fences
		"    %s %s %s\n" + // Line for reversed car parks to still get fences
		"    %s\n", // 80 - set bit 15 to show this is a final return value
		length,
		bytes.GetByte(s.SetID+setIDOffset),
		checkValue,
		bytes.GetWord(255),
		ifTrueValue,
		bytes.GetWord(255), // If the tile is not a station the value of the lower bits is 0xFFFF
		bytes.GetWord(255),
		ifTrueValue,
		bytes.GetWord(FIRST_CAR_PARK_ID),
		bytes.GetWord((FIRST_CAR_PARK_ID + NUM_CAR_PARKS)-1),
		ifTrueValue,
		bytes.GetWord(FIRST_REVERSE_CAR_PARK_ID),
		bytes.GetWord((FIRST_REVERSE_CAR_PARK_ID + NUM_CAR_PARKS)-1),
		ifFalseValue)
}

func (s *StationFenceCallback) GetLines() []string {
	// Fence layouts follow a specific pattern, but we might not be using the one which starts at
	// 0, so we add the offset to the fixed pattern of N/S fences
	result := []string{
		s.getAction("10", bytes.GetCallbackResultByte(s.BaseLayoutOffset+4), bytes.GetCallbackResultByte(s.BaseLayoutOffset+0), 1), // N: true, S: check
		s.getAction("10", bytes.GetCallbackResultByte(s.BaseLayoutOffset+6), bytes.GetCallbackResultByte(s.BaseLayoutOffset+2), 2), // N: false, S: check
		s.getAction("F0", bytes.GetWord(s.SetID+2), bytes.GetWord(s.SetID+1), STATION_FENCE_OFFSET),                                // N: check, S: unknown
	}

	if s.HasDecider {
		result = append(result, GetDecider(s.SetID, s.SetID+3, s.YearCallbackID, s.DefaultSpriteSet))
	}

	return result
}
