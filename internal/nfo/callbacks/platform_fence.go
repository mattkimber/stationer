package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type PlatformFenceCallback struct {
	YearCallbackID int
}

func (pf *PlatformFenceCallback) GetComment() string {
	return "Callback for 4-way platform fences"
}

// Chain:
// Track continuation state -> N / S / Both
// Then decide fences
// Both = callbacks for fence n or !n, then combinations of fences
// N/S = additionally check front or back fence

func (pf *PlatformFenceCallback) getFenceCallback(checkValue, ifTrueValue, ifFalseValue string, callbackID int) string {
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

func (pf *PlatformFenceCallback) GetLines() []string {
	return []string{

		// Callbacks for fence combinations
		// 00 = no fences, 02 = fence to N, 04 = fence to S, 06 = n/s,
		// 08 = fence to rear, 0A = rear/n, 0C = rear/s, 0E = rear/n/s
		// 10 = fence to front,
		//      12 = front/n, 14 = front/s, 16 = front/n/s
		//      18 = front/rear, 1A = front/rear/n, 1C = front/rear/s, 1E = front/rear/n/s
		pf.getFenceCallback("10", "0E 80", "0A 80", 2),                   // F: false, R: true, N: true, S: check
		pf.getFenceCallback("10", "0C 80", "08 80", 3),                   // F: false, R: true, N: false, S: check
		pf.getFenceCallback("F0", bytes.GetWord(2), bytes.GetWord(3), 4), // F: false, R: true, N: check, S: unknown

		pf.getFenceCallback("10", "06 80", "02 80", 5),                   // F: false, R: false, N: true, S: check
		pf.getFenceCallback("10", "04 80", "00 80", 6),                   // F: false, R: false, N: false, S: check
		pf.getFenceCallback("F0", bytes.GetWord(5), bytes.GetWord(6), 7), // F: false, R: false, N: check, S: unknown

		pf.getFenceCallback("0F", bytes.GetWord(4), bytes.GetWord(7), 8), // F: false, R: check, N: unknown, S: unknown

		pf.getFenceCallback("10", "1E 80", "1A 80", 9),                     // F: true, R: true, N: true, S: check
		pf.getFenceCallback("10", "1C 80", "18 80", 10),                    // F: true, R: true, N: false, S: check
		pf.getFenceCallback("F0", bytes.GetWord(9), bytes.GetWord(10), 11), // F: true, R: true, N: check, S: unknown

		pf.getFenceCallback("10", "16 80", "12 80", 12),                     // F: true, R: false, N: true, S: check
		pf.getFenceCallback("10", "14 80", "10 80", 13),                     // F: true, R: false, N: false, S: check
		pf.getFenceCallback("F0", bytes.GetWord(12), bytes.GetWord(13), 14), // F: true, R: false, N: check, S: unknown

		pf.getFenceCallback("0F", bytes.GetWord(11), bytes.GetWord(14), 15), // F: true, R: check, N: unknown, S: unknown

		pf.getFenceCallback("01", bytes.GetWord(15), bytes.GetWord(8), 1), // F: check, R: unknown, N: unknown, S: unknown

		GetDecider(0, 1, pf.YearCallbackID, 0),
	}
}
