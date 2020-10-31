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
	callback += "    04 80 00 00\n"

	// 1 = track continues to S only
	callback += "    00 80 01 01\n"

	// 2 = track continues to N only
	callback += "    02 80 02 02\n"

	// 3 = track continues to N and S
	callback += "    04 80 03 03\n"

	// Default sprite set (shouldn't be used)
	callback += "    04 80"

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
		mtb.getCallback(),
		mtb.getDecider(),
	}
}
