package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type LargeCentralObjectCallback struct {
	SetID            int
	OuterCallbackID  int // the callback ID of the "outer"/southernmost platform
	InnerCallbackID  int // the callback ID of the "inner"/northernmost platform
	MiddleCallbackID int // the callback ID for in-between platforms with objects
	DefaultSpriteSet int
	YearCallbackID   int
	HasDecider       bool
}

func (lco *LargeCentralObjectCallback) GetComment() string {
	return "Callback for multi-tile central platform object"
}

func (lco *LargeCentralObjectCallback) getCallback(positionMask, edgeValue, middleValue string, offset int) string {
	length := 14

	callback := fmt.Sprintf(
		// 81 = byte-sized variable
		// 41 = variable (platform info for this section, counted from northern edge)
		// 08 0F = get position across platform from north (mask)
		// %s = number of ranges other than default
		"* %d 02 04 %s\n"+
			"    81 41 %s\n"+
			"    %s\n",
		length,
		bytes.GetByte(lco.SetID+offset), // ID of this callback - we add the offset to avoid colliding switches
		positionMask,                    // the platform position mask to use
		bytes.GetByte(1),                // number of ranges
	)

	// set ID, low range, high range
	// We return layout 0 for the case where we have a building
	callback += fmt.Sprintf("    %s 00 00 \n", edgeValue)

	// Default sprite set, layout 2 (no building)
	callback += fmt.Sprintf("    %s", middleValue)

	return callback
}

func (lco *LargeCentralObjectCallback) GetLines() []string {
	result := []string{
		lco.getCallback("08 0F", bytes.GetWord(lco.InnerCallbackID), bytes.GetWord(lco.MiddleCallbackID), 1),
		lco.getCallback("08 F0", bytes.GetWord(lco.OuterCallbackID), bytes.GetWord(lco.SetID+1), 2),
	}


	if lco.HasDecider {
		result = append(result, GetDecider(lco.SetID, lco.SetID+2, lco.YearCallbackID, lco.DefaultSpriteSet))
	}

	return result
}
