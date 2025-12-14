package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type MinimumBridgeClearance struct {
	Clearance int
	Layouts   int
}

func (m *MinimumBridgeClearance) GetBytes() int {
	return 2 + m.Layouts
}

func (m *MinimumBridgeClearance) GetString() string {
	output := ""
	for i := 0; i < m.Layouts; i++ {
		output += fmt.Sprintf("%s ", bytes.GetByte(m.Clearance))
	}

	return fmt.Sprintf("20 %s %s // min clearance of %d", bytes.GetByte(m.Layouts), output, m.Clearance)
}
