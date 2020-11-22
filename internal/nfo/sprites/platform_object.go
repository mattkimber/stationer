package sprites

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
)

type PlatformObject struct {
	BaseSpriteID    int
	SpriteFilename  string
	MaxLoadState    int
	InvertDirection bool
	IsStatic        bool
}

func (po *PlatformObject) GetSprite(filename string, num int, swap bool) Sprite {
	xrel := -(SPRITE_WIDTH / 2) - 10

	if swap {
		xrel = 11 - (SPRITE_WIDTH / 2)
	}

	return Sprite{
		Filename: filename,
		X:        SPRITE_WIDTH_WITH_PADDING * num,
		Y:        0,
		Width:    SPRITE_WIDTH,
		Height:   SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     -(SPRITE_HEIGHT / 2) - 1,
	}
}

func (po *PlatformObject) WriteToFile(file *output_file.File, loadState int) {

	filename := fmt.Sprintf("%s_%d_8bpp.png", po.SpriteFilename, loadState)
	if po.IsStatic {
		filename = fmt.Sprintf("%s_8bpp.png", po.SpriteFilename)
	}

	if loadState <= po.MaxLoadState {
		file.AddElement(&Sprites{
			po.GetSprite(filename, 0, po.InvertDirection != true),
			po.GetSprite(filename, 1, po.InvertDirection != false),
		})
	} else {
		file.AddElement(&Blank{Size: 2})
	}

}
