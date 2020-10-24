package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
)

type PylonPlacement struct {
	Bitmask int
}

func (p *PylonPlacement) GetBytes() int {
	return 2
}

func (p *PylonPlacement) GetString() string {
	return fmt.Sprintf("11 %s ", bytes.GetByte(p.Bitmask))
}

type WirePlacement struct {
	Bitmask int
}

func (w *WirePlacement) GetBytes() int {
	return 2
}

func (w *WirePlacement) GetString() string {
	return fmt.Sprintf("14 %s ", bytes.GetByte(w.Bitmask))
}