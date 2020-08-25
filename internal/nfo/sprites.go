package nfo

import (
	"fmt"
)

const (
	FEATURE_STATIONS = 4
)

type Sprites []Sprite

type Sprite struct {
	Filename string
	X int
	Y int
	Width int
	Height int
	XRel int
	YRel int
}

func (s *Sprites) GetLines() []string {
	result := make([]string, 1 + len(*s))
	result[0] = fmt.Sprintf("* 6 01 %s 01 FF %s 00", GetByte(FEATURE_STATIONS), GetByte(len(*s)))


	for idx, spr := range *s {
		result[idx+1] = fmt.Sprintf("%s 8bpp %d %d %d %d %d %d normal chunked", "1x/" + spr.Filename, spr.X, spr.Y,  spr.Width, spr.Height,spr.XRel, spr.YRel)
	}

	return result
}