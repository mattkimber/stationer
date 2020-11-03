package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type AvailabilityYearCallback struct {
	SetID int
	HasDecider bool
	Year int
}

func (ay *AvailabilityYearCallback) GetComment() string {
	return "Callback for availability year"
}

func (ay *AvailabilityYearCallback) getCallback(setID int) string {
	length := 23

	return fmt.Sprintf(
		// 89 = access doubleword
		// 24 = variable (current year zero based)
		// 00 FF FF FF FF = do not adjust the variable
		// 01 = number of ranges other than default (1 in this case)
		"* %d 02 04 %s 89 24 00 FF FF FF FF 01\n"+
			// set ID, low range, high range
			// note that ranges are bytes, as the variable is a byte
			"    00 80 00 00 00 00 %s \n"+ // 0 = station not available
			"    01 80",                   // 1 = station available
		length,
		bytes.GetByte(setID), // The callback decider is given the SetID
		bytes.GetDouble(ay.Year - 1),   // Last year station is not available
	)
}

func (ay *AvailabilityYearCallback) GetLines() []string {
	if ay.HasDecider {
		return []string{
			ay.getCallback(ay.SetID+1),
			GetDecider(ay.SetID, 0, ay.SetID+1, 0),
		}
	}

	// No decider, piggyback one from another callback
	return []string{
		ay.getCallback(ay.SetID),
	}
}
