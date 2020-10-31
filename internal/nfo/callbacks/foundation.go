package callbacks

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type FoundationCallback struct {
	SetID            int
}

func (f *FoundationCallback) GetComment() string {
	return "Callback for custom foundation"
}


func (f *FoundationCallback) getCallback() string {
	length := 14

	return fmt.Sprintf(
		// 81 = access lowest byte
		// 10 = variable (foundation info)
		// 00 FF = do not adjust the variable
		// 01 = number of ranges other than default (1 in this case)
		"* %d 02 04 %s 81 10 00 FF 01\n"+
			// set ID, low range, high range
			// note that ranges are bytes, as the variable is a byte
			"    %s 02 02 \n"+
			"    %s",
		length,
		bytes.GetByte(f.SetID), // The callback decider is given the SetID
		bytes.GetWord(1), // 1 = the foundations
		bytes.GetWord(0), // 0 = the base sprite set
	)
}



func (f *FoundationCallback) GetLines() []string {
	return []string{
		f.getCallback(),
	}
}
