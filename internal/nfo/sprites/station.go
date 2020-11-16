package sprites

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
)

type StationSprites struct {
	ID                    int
	SpriteFilename        string
	MaxLoadState          int
	HasFences             bool
}

func (s *StationSprites) GetSprite(filename string, num int, swap bool) Sprite {
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

func (s *StationSprites) WriteToFile(file *output_file.File) {
	platformSprites := 4
	if s.HasFences {
		platformSprites = 8
	}

	file.AddElement(&Spritesets{ID: s.ID, NumSets: s.MaxLoadState + 1, NumSprites: platformSprites})

	for i := 0; i <= s.MaxLoadState; i++ {
		filename := fmt.Sprintf("%s_%d_8bpp.png", s.SpriteFilename, i)

		// Non-fence sprites
		file.AddElement(&Sprites{
			s.GetSprite(filename, 0, false),
			s.GetSprite(filename, 1, false),
			s.GetSprite(filename, 2, true),
			s.GetSprite(filename, 3, true),
		})

		// Fence sprites
		if s.HasFences {
			fenceFilename := fmt.Sprintf("%s_%d_fence_8bpp.png", s.SpriteFilename, i)
			file.AddElement(&Sprites{
				s.GetSprite(fenceFilename, 0, false),
				s.GetSprite(fenceFilename, 1, false),
				s.GetSprite(fenceFilename, 2, true),
				s.GetSprite(fenceFilename, 3, true),
			})
		}
	}
}