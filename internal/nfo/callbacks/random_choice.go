package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type RandomChoiceCallback struct {
	SetID int
	// ResultIDs must have length of a power of 2
	ResultIDs        []int
	HasDecider       bool
	YearCallbackID   int
	DefaultSpriteSet int
}

func (rcb *RandomChoiceCallback) GetComment() string {
	return "Random choice callback"
}

func (rcb *RandomChoiceCallback) getCallback() string {
	length := 7 + (len(rcb.ResultIDs) * 2)

	callback := fmt.Sprintf(
		// 80 = random choice based on own variables

		// 00 = count (byte), only used for vehicles so not relevant here
		// 00 = no random triggers (byte)
		// 10 = random bits to use (byte)
		// %s = number of options (must be a power of 2)

		// %s = number of ranges other than default
		"* %d 02 04 %s\n"+
			"    80 00\n"+
			"    10\n"+
			"    %s\n",
		length,
		bytes.GetByte(rcb.SetID),          // ID of this callback
		bytes.GetByte(len(rcb.ResultIDs)), // number of ranges
	)

	for _, resultID := range rcb.ResultIDs {
		// Add the set IDs to choose from randomly
		callback += fmt.Sprintf("    %s\n", bytes.GetWord(resultID))
	}

	return callback
}

func (rcb *RandomChoiceCallback) GetLines() []string {
	result := []string{
		rcb.getCallback(),
	}

	if rcb.HasDecider {
		result = append(result, GetDecider(rcb.SetID+1, rcb.SetID, rcb.YearCallbackID, rcb.DefaultSpriteSet))
	}

	return result
}
