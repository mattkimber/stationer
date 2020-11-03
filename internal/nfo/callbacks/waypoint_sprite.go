package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type WaypointSpriteCallback struct {
	YearCallbackID int
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

func (wsc *WaypointSpriteCallback) GetLines() []string {
	return []string{
		wsc.getCallback(),
		GetDecider(2, 1, wsc.YearCallbackID, 0),
	}
}
