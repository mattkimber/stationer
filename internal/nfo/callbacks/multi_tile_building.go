package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type MultiTileBuildingCallback struct {
	SetID  int
	Length int
}

func (mtb *MultiTileBuildingCallback) GetComment() string {
	return "Callback for multi-tile building"
}

func (mtb *MultiTileBuildingCallback) getCallback() string {
	length := 10 + ((mtb.Length - 1) * 4)

	callback := fmt.Sprintf(
		// 89 = doubleword variable
		// 41 = variable (platform info for this section, counted from northern edge)
		// 00 0F = get position along platform from north (mask)
		// %s = number of ranges other than default
		"* %d 02 04 %s\n"+
			"    81 41 00 0F\n"+
			"    %s\n",
		length,
		bytes.GetByte(mtb.SetID+1), // We add 1 to the decider ID to give the inner callback ID
		bytes.GetByte(mtb.Length-1),
	)

	for i := 1; i < mtb.Length; i++ {
		callback += fmt.Sprintf(
			// set ID, low range, high range
			// note that ranges are doubles, as the variable is a double
			"    %s 80 %s %s \n",
			bytes.GetByte(i*2), // Set ID of the building tile
			bytes.GetByte(i),
			bytes.GetByte(i),
		)
	}

	// Default sprite set = 0
	callback += "    00 80"

	return callback
}

func (mtb *MultiTileBuildingCallback) getDecider() string {
	// Callbacks require a callback decider. This will be passed the type
	// of callback (station layout = 14) and be responsible for routing it
	// to the correct callback.
	length := 17

	return fmt.Sprintf(
		"* %d 02 04 %s 85 0C 00 FF FF 01\n"+
			"    %s 00 14 00 14 00\n"+
			"    00 00", // Return the default sprite if we don't trigger any callback
		length,
		bytes.GetByte(mtb.SetID),   // The callback decider is given the SetID
		bytes.GetByte(mtb.SetID+1), // The actual callback is SetID + 1
	)
}

func (mtb *MultiTileBuildingCallback) GetLines() []string {
	return []string{
		mtb.getCallback(),
		mtb.getDecider(),
	}
}
