package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

const (
	STATION_FENCE_OFFSET      = 3
	FIRST_CAR_PARK_ID         = 157
	FIRST_REVERSE_CAR_PARK_ID = 172
	NUM_CAR_PARKS             = 4

	FIRST_WOODEN_INNER_PLATFORM = 14
	FIRST_WOODEN_OUTER_PLATFORM = 23

	FIRST_CONCRETE_INNER_PLATFORM = 55
	FIRST_CONCRETE_OUTER_PLATFORM = 72

	FIRST_MODERN_INNER_PLATFORM = 111
)

type StationFenceCallback struct {
	SetID                   int
	BaseLayoutOffset        int
	DefaultSpriteSet        int
	YearCallbackID          int
	HasDecider              bool
	UseRailPresenceForNorth bool
	UseRailPresenceForSouth bool
}

func (s *StationFenceCallback) GetComment() string {
	return "Callback for station fence choice"
}

func (s *StationFenceCallback) getRailPresenceAction(mask, ifTrueValue, ifFalseValue string, setIDOffset int) string {
	length := 17

	// 85 = access lowest word of station variable
	// 45 = rail continuation info of nearby tile
	return fmt.Sprintf("* %d 02 04 %s\n"+
		"    85 45\n"+ // Check variable 45
		"    00\n"+ // No variable adjustment
		"    00 %s\n"+ // mask higher bits
		"    01\n"+ // 1 non-default option
		"    %s %s %s\n"+ // Everything is FF
		"    %s\n", // 80 - set bit 15 to show this is a final return value
		length,
		bytes.GetByte(s.SetID+setIDOffset),
		mask,
		ifFalseValue,
		bytes.GetWord(1), // Any non-zero value indicates track continuation
		bytes.GetWord(65535),
		ifTrueValue)
}

func (s *StationFenceCallback) getStationPresenceAction(checkValue, ifTrueValue, ifFalseValue string, setIDOffset int) string {
	length := 48

	// North should show fences against inner platforms
	woodenPlatformExclusionStart := FIRST_WOODEN_INNER_PLATFORM
	woodenPlatformExclusionCount := FIRST_WOODEN_OUTER_PLATFORM - FIRST_WOODEN_INNER_PLATFORM - 1

	concretePlatformExclusionStart := FIRST_CONCRETE_INNER_PLATFORM
	modernPlatformExclusionStart := FIRST_MODERN_INNER_PLATFORM
	platformExclusionCount := FIRST_CONCRETE_OUTER_PLATFORM - FIRST_CONCRETE_INNER_PLATFORM - 1

	if checkValue == "10" {
		// South should show fences against outer platforms
		woodenPlatformExclusionStart = FIRST_WOODEN_OUTER_PLATFORM
		concretePlatformExclusionStart = FIRST_CONCRETE_OUTER_PLATFORM
		modernPlatformExclusionStart = FIRST_MODERN_INNER_PLATFORM + platformExclusionCount

		// No hut or billboard for outer platforms
		woodenPlatformExclusionCount = woodenPlatformExclusionCount - 2
		// + No cafe or planter for modern outer platforms
		platformExclusionCount = platformExclusionCount - 4
	}

	// 85 = access lowest word of station variable
	// 68 = station info of nearby tile
	return fmt.Sprintf("* %d 02 04 %s\n"+
		"    85 68 %s\n"+ // Check variable 68 for tile to -1/0
		"    00 %s\n"+ // mask bits
		"    06\n"+ // 6 non-default options
		"    %s %s %s\n"+ // Everything is FF
		"    %s %s %s\n"+ // Line for car parks to still get fences
		"    %s %s %s\n"+ // Line for wooden inner/outer to still get fences
		"    %s %s %s\n"+ // Line for concrete inner/outer to still get fences
		"    %s %s %s\n"+ // Line for modern inner/outer to still get fences
		"    %s %s %s\n"+ // Line for reversed car parks to still get fences
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
		bytes.GetWord((FIRST_CAR_PARK_ID+NUM_CAR_PARKS)-1),
		ifTrueValue,
		bytes.GetWord(woodenPlatformExclusionStart),
		bytes.GetWord(woodenPlatformExclusionStart+woodenPlatformExclusionCount),
		ifTrueValue,
		bytes.GetWord(concretePlatformExclusionStart),
		bytes.GetWord(concretePlatformExclusionStart+platformExclusionCount),
		ifTrueValue,
		bytes.GetWord(modernPlatformExclusionStart),
		bytes.GetWord(modernPlatformExclusionStart+platformExclusionCount),
		ifTrueValue,
		bytes.GetWord(FIRST_REVERSE_CAR_PARK_ID),
		bytes.GetWord((FIRST_REVERSE_CAR_PARK_ID+NUM_CAR_PARKS)-1),
		ifFalseValue)
}

func (s *StationFenceCallback) GetLines() []string {
	// Fence layouts follow a specific pattern, but we might not be using the one which starts at
	// 0, so we add the offset to the fixed pattern of N/S fences
	result := make([]string, 3)

	if s.UseRailPresenceForSouth {
		result[0] = s.getRailPresenceAction("04", bytes.GetCallbackResultByte(s.BaseLayoutOffset+4), bytes.GetCallbackResultByte(s.BaseLayoutOffset+0), 1) // N: true, S: check
		result[1] = s.getRailPresenceAction("04", bytes.GetCallbackResultByte(s.BaseLayoutOffset+6), bytes.GetCallbackResultByte(s.BaseLayoutOffset+2), 2) // N: false, S: check
	} else {
		result[0] = s.getStationPresenceAction("10", bytes.GetCallbackResultByte(s.BaseLayoutOffset+4), bytes.GetCallbackResultByte(s.BaseLayoutOffset+0), 1) // N: true, S: check
		result[1] = s.getStationPresenceAction("10", bytes.GetCallbackResultByte(s.BaseLayoutOffset+6), bytes.GetCallbackResultByte(s.BaseLayoutOffset+2), 2) // N: false, S: check
	}

	if s.UseRailPresenceForNorth {
		result[2] = s.getRailPresenceAction("08", bytes.GetWord(s.SetID+2), bytes.GetWord(s.SetID+1), STATION_FENCE_OFFSET) // N: check, S: unknown
	} else {
		result[2] = s.getStationPresenceAction("F0", bytes.GetWord(s.SetID+2), bytes.GetWord(s.SetID+1), STATION_FENCE_OFFSET) // N: check, S: unknown
	}

	if s.HasDecider {
		result = append(result, GetDecider(s.SetID, s.SetID+3, s.YearCallbackID, s.DefaultSpriteSet))
	}

	return result
}
