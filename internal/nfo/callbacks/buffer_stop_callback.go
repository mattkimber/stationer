package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type BufferStopCallback struct {}

func (mtb *BufferStopCallback) GetComment() string {
	return "Callback for buffer stop"
}

// Chain:
// Track continuation state -> N / S / Both
// Then decide fences
// Both = callbacks for fence n or !n, then combinations of fences
// N/S = additionally check front or back fence


func (s *BufferStopCallback) getFenceCallback(checkValue, ifTrueValue, ifFalseValue string, callbackID int) string {
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
		bytes.GetByte(callbackID),
		checkValue,
		bytes.GetWord(65535),
		ifTrueValue,
		bytes.GetWord(65535), // If the tile is not a station the value of the lower bits is 0xFFFF
		bytes.GetWord(65535),
		ifFalseValue)
}

func (mtb *BufferStopCallback) getCallback() string {
	length := 10 + (4*4)

	callback := fmt.Sprintf(
		// 81 = get lowest byte of variable
		// 45 = variable (track continuation information)
		// 00 03 = 00000011 - platform continuation +/- length
		// 04 = 4 ranges (other than default)
		"* %d 02 04 %s\n" +
			"    81 45 00 03\n" +
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
	callback += "    00 80 01 01\n"

	// 2 = track continues to N only
	// 0E = 14 (as we have 7 fence sprite combinations, including no fence, for 2 directions, for 1 type of buffer stop)
	callback += "    0E 80 02 02\n"

	// 3 = track continues to N and S
	// 1C = 28 (as we have 7 fence sprite combinations, including no fence, for 2 directions, for 2 types of buffer stop)
	callback += "    0A 00 03 03\n"

	// Default sprite set (shouldn't be used)
	callback += "    00 80"

	return callback
}


func (mtb *BufferStopCallback) getDecider() string {
	// Callbacks require a callback decider. This will be passed the type
	// of callback (station layout = 14) and be responsible for routing it
	// to the correct callback.
	length := 17

	return fmt.Sprintf(
		"* %d 02 04 %s 85 0C 00 FF FF 01\n"+
			"    %s 00 14 00 14 00\n"+
			"    00 00", // Return the default sprite if we don't trigger any callback
		length,
		bytes.GetByte(0), // The callback decider is given this ID
		bytes.GetByte(1), // The first callback in the chain is this ID
	)
}

func (mtb *BufferStopCallback) GetLines() []string {
	return []string{
		// Callbacks for "both"
		// 1C = no fences, 1E = fence to N, 20 = fence to S, 22 = fence to N/S
		mtb.getFenceCallback("10", "22 80", "1E 80", 11), // N: true, S: check
		mtb.getFenceCallback("10", "20 80", "1C 80", 12), // N: false, S: check
		mtb.getFenceCallback("F0", bytes.GetWord(11), bytes.GetWord(12), 10), // N: check, S: unknown
		mtb.getCallback(),
		mtb.getDecider(),
	}
}
