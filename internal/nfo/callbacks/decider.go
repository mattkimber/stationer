package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

func GetDecider(deciderID, layoutCallbackID, yearCallbackID, defaultSpriteset int) string {
	// Callbacks require a callback decider. This will be passed the type
	// of callback (station layout = 14) and be responsible for routing it
	// to the correct callback.
	length := 11
	ranges := 0

	yearCallback, layoutCallback := "", ""

	if yearCallbackID != 0 {
		yearCallback = fmt.Sprintf("    %s 00 13 00 13 00\n", bytes.GetByte(yearCallbackID)) // var 13 = availability of station
		length += 6
		ranges += 1
	}

	if layoutCallbackID != 0 {
		layoutCallback = fmt.Sprintf("    %s 00 14 00 14 00\n", bytes.GetByte(layoutCallbackID))
		length += 6
		ranges += 1
	}

	return fmt.Sprintf(
		"* %d 02 04 %s 85 0C 00 FF FF %s\n"+
			yearCallback +
			layoutCallback +
			"    %s 00", // Return the default sprite if we don't trigger any callback
		length,
		bytes.GetByte(deciderID),   // The callback decider is given the SetID
		bytes.GetByte(ranges),      // Number of ranges can change if we have more than one callback
									// If this isn't right you will get weird results (default sprites, etc.)
		bytes.GetByte(defaultSpriteset), // If you return 0 for platforms rather than the correct spriteset, you
		                                 // get incorrect sprites in the purchase menu
	)
}
