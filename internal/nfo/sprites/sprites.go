package sprites

import (
	"fmt"
)

const (
	FEATURE_STATIONS = 4
)

type Sprites []Sprite

type Sprite struct {
	Filename string
	X        int
	Y        int
	Width    int
	Height   int
	XRel     int
	YRel     int
}

func (s *Sprites) GetComment() string {
	return "Sprites"
}

func (s *Sprites) GetLines() []string {
	result := make([]string, len(*s))

	for idx, spr := range *s {
		if spr.Filename != "" {
			zi1 := fmt.Sprintf("%s 8bpp %d %d %d %d %d %d normal chunked", "1x/"+spr.Filename, spr.X, spr.Y, spr.Width, spr.Height, spr.XRel, spr.YRel)
			zi2 := fmt.Sprintf("| %s 8bpp %d %d %d %d %d %d zi2 chunked", "2x/"+spr.Filename, spr.X*2, spr.Y*2, spr.Width*2, spr.Height*2, spr.XRel*2, (spr.YRel*2)-1)

			result[idx] = zi1 + "\n" + zi2
		} else {
			result[idx] = "files/transparent.png 8bpp 0 0 1 1 0 0 normal chunked"
		}
	}

	return result
}
