package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

// Little/lots threshold
type LittleLotsThreshold struct {
	Amount int
}

func (l *LittleLotsThreshold) GetBytes() int {
	return 3
}

func (l *LittleLotsThreshold) GetString() string {
	return fmt.Sprintf("10 %s ", bytes.GetWord(l.Amount))
}
