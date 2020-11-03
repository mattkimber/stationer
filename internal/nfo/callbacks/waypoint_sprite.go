package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type WaypointSpriteCallback struct {
}

func (wsc *WaypointSpriteCallback) GetComment() string {
	return "Callback for multi-tile building"
}

func (wsc *WaypointSpriteCallback) getCallback() string {
	length := 14

	callback := fmt.Sprintf(
		// 81 = byte-sized variable
		// 41 = variable (platform info for this section, counted from northern edge)
		// 08 0F = get position across platform from north (mask)
		// %s = number of ranges other than default
		"* %d 02 04 %s\n"+
			"    81 41 08 0F\n"+
			"    %s\n",
		length,
		bytes.GetByte(1), // ID of this callback
		bytes.GetByte(1), // number of ranges
	)

	// set ID, low range, high range
	// We return layout 0 for the case where we have a building
	callback += "    00 80 00 00 \n"

	// Default sprite set, layout 2 (no building)
	callback += "    02 80"

	return callback
}

func (wsc *WaypointSpriteCallback) getDecider() string {
	// Callbacks require a callback decider. This will be passed the type
	// of callback (station layout = 14) and be responsible for routing it
	// to the correct callback.
	length := 17

	return fmt.Sprintf(
		"* %d 02 04 %s 85 0C 00 FF FF 01\n"+
			"    %s 00 14 00 14 00\n"+
			"    00 00", // Return the default sprite if we don't trigger any callback
		length,
		bytes.GetByte(2), // The callback decider is given an ID of 2
		bytes.GetByte(1), // The actual callback is SetID + 1
	)
}

func (wsc *WaypointSpriteCallback) GetLines() []string {
	return []string{
		wsc.getCallback(),
		wsc.getDecider(),
	}
}
