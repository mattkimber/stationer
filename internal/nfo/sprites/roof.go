package sprites

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
)

type StationRoof struct {
	ID             int
	SpriteFilename string
	MaxLoadState   int
	RoofType       string
}

func (s *StationRoof) GetSprite(filename string, num int, swap bool) Sprite {
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

func (s *StationRoof) GetRoofSprite(filename string, num int) Sprite {
	xrel := 1 - (BUILDING_SPRITE_WIDTH / 2)
	yrel := -(BUILDING_SPRITE_HEIGHT / 2) - 6

	return Sprite{
		Filename: filename,
		X:        BUILDING_SPRITE_WIDTH_WITH_PADDING * num,
		Y:        0,
		Width:    BUILDING_SPRITE_WIDTH,
		Height:   BUILDING_SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     yrel,
	}
}

func (s *StationRoof) WriteToFile(file *output_file.File) {

	file.AddElement(&Spritesets{ID: 0, NumSets: s.MaxLoadState + 1, NumSprites: 4 + (3 * 4)})

	for i := 0; i <= s.MaxLoadState; i++ {

		// Sprites with roof support
		filename := fmt.Sprintf("%s_%d_roof_8bpp.png", s.SpriteFilename, i)

		file.AddElement(&Sprites{
			s.GetSprite(filename, 0, false),
			s.GetSprite(filename, 1, false),
			s.GetSprite(filename, 2, true),
			s.GetSprite(filename, 3, true),
		})

		roofSubtypes := []string{"single", "double_n", "double_s"}

		for _, subtype := range roofSubtypes {
			// 3 roof sprites (single, N, S) - repeated as we have multiple load states
			filename = fmt.Sprintf("roof_%s_%s_8bpp.png", s.RoofType, subtype)
			file.AddElement(&Sprites{
				s.GetRoofSprite(filename, 0),
				s.GetRoofSprite(filename, 1),
			})

			// And again for the glass
			filename = fmt.Sprintf("roof_%s_%s_glass_8bpp.png", s.RoofType, subtype)
			file.AddElement(&Sprites{
				s.GetRoofSprite(filename, 0),
				s.GetRoofSprite(filename, 1),
			})
		}
	}
}
