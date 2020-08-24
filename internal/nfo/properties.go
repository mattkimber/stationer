package nfo

import "fmt"

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
}

type SpriteDirection struct {
	GroundSprite int
	Foreground BoundingBox
	Background BoundingBox
}

type SpriteLayout struct {
	EastWest SpriteDirection
	NorthSouth SpriteDirection
}

func (sl *SpriteLayout) GetBytes() int {
	return (BOUNDING_BOX_BYTES * 4) +
		(TILE_DIRECTION_END_BYTES * 2) +
		(GROUNDSPRITE_BYTES * 2) + 2
}

func (bb *BoundingBox) GetString(coda int) string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s 00 00",
		GetByte(bb.XOffset),
		GetByte(bb.YOffset),
		GetByte(bb.ZOffset),
		GetByte(bb.X),
		GetByte(bb.Y),
		GetByte(bb.Z),
		GetByte(coda),
		GetByte(132),
	)
}

func (sd *SpriteDirection) GetString(bg_coda int, fg_coda int) string {
	return fmt.Sprintf("%s %s %s 80",
		GetWord(sd.GroundSprite),
		sd.Background.GetString(bg_coda),
		sd.Foreground.GetString(fg_coda),
	)
}

func (sl *SpriteLayout) GetString() string {
	return fmt.Sprintf("09 02 %s %s",
		sl.EastWest.GetString(45, 46),
		sl.NorthSouth.GetString(47, 48),
		)
}