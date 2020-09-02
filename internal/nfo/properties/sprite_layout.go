package properties

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/bytes"
	"strings"
)

type BoundingBox struct {
	XOffset int
	YOffset int
	ZOffset int
	X int
	Y int
	Z int
	SpriteNumber int
}

func (bb *BoundingBox) GetString() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s ",
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
	Sprites []BoundingBox
}

type SpriteLayout struct {
	EastWest   SpriteDirection
	NorthSouth SpriteDirection
}

func (sl *SpriteLayout) GetBytes() int {
	boundingBoxes := len(sl.EastWest.Sprites) + len(sl.NorthSouth.Sprites)

	return (BOUNDING_BOX_BYTES * boundingBoxes) +
		(TILE_DIRECTION_END_BYTES * 2) +
		(GROUNDSPRITE_BYTES * 2) + 2
}

func (sd *SpriteDirection) GetString() string {
	sb := strings.Builder{}
	for _, bb := range sd.Sprites {
		sb.WriteString(bb.GetString())
	}

	return fmt.Sprintf("%s %s80 ",
		bytes.GetDouble(sd.GroundSprite),
		sb.String(),
	)
}

func (sl *SpriteLayout) GetString() string {
	return fmt.Sprintf("09 02 %s%s",
		sl.EastWest.GetString(),
		sl.NorthSouth.GetString(),
	)
}