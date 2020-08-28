package nfo

import (
	"fmt"
	"strings"
)

const (
	BOUNDING_BOX_BYTES = 10
	TILE_DIRECTION_END_BYTES = 1
	GROUNDSPRITE_BYTES = 4
)

type Property interface {
	GetBytes() int
	GetString() string
}

type ClassID struct {
	ID string
}

func (c *ClassID) GetBytes() int {
	return 1 + len(c.ID)
}

func (c *ClassID) GetString() string {
	return fmt.Sprintf("08 \"%s\" ", c.ID)
}

type BoundingBox struct {
	XOffset int
	YOffset int
	ZOffset int
	X int
	Y int
	Z int
	SpriteNumber int
}

type SpriteDirection struct {
	GroundSprite int
	Sprites []BoundingBox
}

type SpriteLayout struct {
	EastWest SpriteDirection
	NorthSouth SpriteDirection
}

func (sl *SpriteLayout) GetBytes() int {
	boundingBoxes := len(sl.EastWest.Sprites) + len(sl.NorthSouth.Sprites)

	return (BOUNDING_BOX_BYTES * boundingBoxes) +
		(TILE_DIRECTION_END_BYTES * 2) +
		(GROUNDSPRITE_BYTES * 2) + 2
}

func (bb *BoundingBox) GetString() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s",
		GetByte(bb.XOffset),
		GetByte(bb.YOffset),
		GetByte(bb.ZOffset),
		GetByte(bb.X),
		GetByte(bb.Y),
		GetByte(bb.Z),
		GetWord(bb.SpriteNumber),
	)
}

func (sd *SpriteDirection) GetString() string {
	sb := strings.Builder{}
	for _, bb := range sd.Sprites {
		sb.WriteString(bb.GetString())
	}

	return fmt.Sprintf("%s %s 80",
		GetWord(sd.GroundSprite),
		sb.String(),
	)
}

func (sl *SpriteLayout) GetString() string {
	return fmt.Sprintf("09 02 %s %s",
		sl.EastWest.GetString(),
		sl.NorthSouth.GetString(),
		)
}