package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type BufferStopCallback struct {
	YearCallbackID          int
	UseRailPresenceForNorth bool
	UseRailPresenceForSouth bool
}

func (mtb *BufferStopCallback) GetComment() string {
	return "Callback for buffer stop"
}

// Chain:
// Track continuation state -> N / S / Both
// Then decide fences
// Both = callbacks for fence n or !n, then combinations of fences
// N/S = additionally check front or back fence

func (mtb *BufferStopCallback) getRailPresenceAction(mask, ifTrueValue, ifFalseValue string, callbackID int) string {
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
		bytes.GetByte(callbackID),
		mask,
		ifFalseValue,
		bytes.GetWord(1), // Any non-zero value indicates track continuation
		bytes.GetWord(65535),
		ifTrueValue)
}

func (mtb *BufferStopCallback) getStationPresenceAction(checkValue, ifTrueValue, ifFalseValue string, callbackID int) string {
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
		bytes.GetByte(callbackID),
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

func (mtb *BufferStopCallback) getCallback() string {
	length := 10 + (4 * 4)

	callback := fmt.Sprintf(
		// 81 = get lowest byte of variable
		// 45 = variable (track continuation information)
		// 00 03 = 00000011 - platform continuation +/- length
		// 04 = 4 ranges (other than default)
		"* %d 02 04 %s\n"+
			"    81 45 00 03\n"+
			"    04\n",
		length,
		bytes.GetByte(1), // This is ID 1
	)

	// First option is shown in purchase menu
	// so we explicitly return the no-tracks case
	// which directs to the both-sides fence decider
	// callback
	//
	// (Note this is an intermediate result so we have 00 not 80 in the high byte)
	callback += "    0A 00 00 00\n"

	// 1 = track continues to S only
	callback += "    1E 00 01 01\n"

	// 2 = track continues to N only
	// 0E = 14 (as we have 7 fence sprite combinations, including no fence, for 2 directions, for 1 type of buffer stop)
	callback += "    14 00 02 02\n"

	// 3 = track continues to N and S
	// 1C = 28 (as we have 7 fence sprite combinations, including no fence, for 2 directions, for 2 types of buffer stop)
	callback += "    0A 00 03 03\n"

	// Default sprite set (shouldn't be used)
	callback += "    00 80"

	return callback
}

func (mtb *BufferStopCallback) getBytes(offset int) string {
	return bytes.GetCallbackResultByte(offset)
}

func (mtb *BufferStopCallback) GetLines() []string {
	northAction := mtb.getStationPresenceAction
	southAction := mtb.getStationPresenceAction
	checkValueN, checkValueS := "F0", "10"

	if mtb.UseRailPresenceForNorth {
		northAction = mtb.getRailPresenceAction
		checkValueN = "08"
	}

	if mtb.UseRailPresenceForSouth {
		southAction = mtb.getRailPresenceAction
		checkValueS = "04"
	}

	return []string{
		// Callbacks for "both"
		// 0x20 = no fences, 0x22 = fence to N, 0x24 = fence to S, 0x26 = fence to N/S
		southAction(checkValueS, mtb.getBytes(38), mtb.getBytes(34), 11),                     // N: true, S: check
		southAction(checkValueS, mtb.getBytes(36), mtb.getBytes(32), 12),                     // N: false, S: check
		northAction(checkValueN, bytes.GetWord(11), bytes.GetWord(12), 10), // N: check, S: unknown

		// Callbacks for "N"
		// 0x10 = no fences, 0x12 = fence to N, 0x14 = fence to S, 0x16 = n/s, 0x18 = fence to rear, 0x1A = rear/n, 0x1C = rear/s, 0x1E = rear/n/s
		southAction(checkValueS, mtb.getBytes(30), mtb.getBytes(26), 21),                     // R: true, N: true, S: check
		southAction(checkValueS, mtb.getBytes(28), mtb.getBytes(24), 22),                     // R: true, N: false, S: check
		northAction(checkValueN, bytes.GetWord(21), bytes.GetWord(22), 23), // R: true, N: check, S: unknown

		southAction(checkValueS, mtb.getBytes(22), mtb.getBytes(18), 24),                     // R: false, N: true, S: check
		southAction(checkValueS, mtb.getBytes(20), mtb.getBytes(16), 25),                     // R: false, N: false, S: check
		northAction(checkValueN, bytes.GetWord(24), bytes.GetWord(25), 26), // R: false, N: check, S: unknown

		mtb.getStationPresenceAction("01", bytes.GetWord(23), bytes.GetWord(26), 20), // R: check, N: unknown, S: unknown

		// Callbacks for "S"
		// 0x00 = no fences, 0x02 = fence to N, 0x04 = fence to S, 0x06 = n/s, 0x08 = fence to rear, 0x0A = rear/n, 0x0C = rear/s, 0x0E = rear/n/s
		southAction(checkValueS, mtb.getBytes(14), mtb.getBytes(10), 31),                     // R: true, N: true, S: check
		southAction(checkValueS, mtb.getBytes(12), mtb.getBytes(8), 32),                     // R: true, N: false, S: check
		northAction(checkValueN, bytes.GetWord(31), bytes.GetWord(32), 33), // R: true, N: check, S: unknown

		southAction(checkValueS, mtb.getBytes(6), mtb.getBytes(2), 34),                     // R: false, N: true, S: check
		southAction(checkValueS, mtb.getBytes(4), mtb.getBytes(0), 35),                     // R: false, N: false, S: check
		northAction(checkValueN, bytes.GetWord(34), bytes.GetWord(35), 36), // R: false, N: check, S: unknown

		mtb.getStationPresenceAction("0F", bytes.GetWord(33), bytes.GetWord(36), 30), // R: check, N: unknown, S: unknown

		mtb.getCallback(),
		GetDecider(0, 1, mtb.YearCallbackID, 0),
	}
}
