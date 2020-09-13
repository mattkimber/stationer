package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
	"strings"
)

type BoundingBox struct {
	XOffset      int
	YOffset      int
	ZOffset      int
	X            int
	Y            int
	Z            int
	SpriteNumber int
}

func (bb *BoundingBox) GetString() string {
	return fmt.Sprintf("\n      %s %s %s %s %s %s %s ",
		bytes.GetByte(bb.XOffset),
		bytes.GetByte(bb.YOffset),
		bytes.GetByte(bb.ZOffset),
		bytes.GetByte(bb.X),
		bytes.GetByte(bb.Y),
		bytes.GetByte(bb.Z),
		bytes.GetDouble(bb.SpriteNumber),
	)
}

type SpriteDirection struct {
	GroundSprite int
	Sprites      []BoundingBox
}

type LayoutEntry struct {
	EastWest   SpriteDirection
	NorthSouth SpriteDirection
}

type SpriteLayout struct {
	Entries []LayoutEntry
}

func (sl *SpriteLayout) GetBytes() int {
	bytes := 2 // Bytes for the action 9 and num-entries
	for _, entry := range sl.Entries {
		bytes += entry.GetBytes()
	}
	return bytes
}

func (sl *SpriteLayout) GetString() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("09 %s ", bytes.GetByte(len(sl.Entries)*2)))

	for _, entry := range sl.Entries {
		sb.WriteString(entry.GetString())
	}

	return sb.String()
}

func (le *LayoutEntry) GetBytes() int {
	boundingBoxes := len(le.EastWest.Sprites) + len(le.NorthSouth.Sprites)

	return (BOUNDING_BOX_BYTES * boundingBoxes) +
		(TILE_DIRECTION_END_BYTES * 2) +
		(GROUNDSPRITE_BYTES * 2)
}

func (sd *SpriteDirection) GetString() string {
	sb := strings.Builder{}
	for _, bb := range sd.Sprites {
		sb.WriteString(bb.GetString())
	}

	return fmt.Sprintf("\n      %s %s\n      80 ",
		bytes.GetDouble(sd.GroundSprite),
		sb.String(),
	)
}

func (le *LayoutEntry) GetString() string {
	return fmt.Sprintf("%s%s",
		le.EastWest.GetString(),
		le.NorthSouth.GetString(),
	)
}
