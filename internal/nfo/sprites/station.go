package sprites

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
)

type StationSprite struct {
	Filename string
	HasFences bool
	MaxLoadState int
}

type StationSprites struct {
	ID                    int
	BaseFilename          string
	Sprites 			  []StationSprite
	SpriteMap             map[string]int // Will be populated after the sprites have been written to a file
	MaxLoadState          int
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

func (s *StationSprites) getTotalSprites() int {
	total := 0

	for _, sprite := range s.Sprites {
		total += 4
		if sprite.HasFences {
			// Add another 4 for the fences
			total += 4
		}
	}

	return total
}

func (s *StationSprites) WriteToFile(file *output_file.File) {
	platformSprites := s.getTotalSprites()
	s.SpriteMap = make(map[string]int)

	file.AddElement(&Spritesets{ID: s.ID, NumSets: s.MaxLoadState + 1, NumSprites: platformSprites})


	for i := 0; i <= s.MaxLoadState; i++ {
		sprites := 0

		for _, spr := range s.Sprites {
			// Populate the map of where sprites begin and end
			s.SpriteMap[spr.Filename] = sprites

			filename := fmt.Sprintf("%s_%s_%d_8bpp.png", s.BaseFilename, spr.Filename, i)

			if i <= spr.MaxLoadState {
				// Non-fence sprites
				file.AddElement(&Sprites{
					s.GetSprite(filename, 0, false),
					s.GetSprite(filename, 1, false),
					s.GetSprite(filename, 2, true),
					s.GetSprite(filename, 3, true),
				})
			} else {
				// Add blank pseudosprites
				file.AddElement(&Blank{Size: 4})
			}

			sprites += 4

			// Fence sprites
			if spr.HasFences {
				fenceFilename := fmt.Sprintf("%s_%s_%d_fence_8bpp.png", s.BaseFilename, spr.Filename, i)
				if i <= spr.MaxLoadState {
					file.AddElement(&Sprites{
						s.GetSprite(fenceFilename, 0, false),
						s.GetSprite(fenceFilename, 1, false),
						s.GetSprite(fenceFilename, 2, true),
						s.GetSprite(fenceFilename, 3, true),
					})
				} else {
					// Add blank pseudosprites
					file.AddElement(&Blank{Size: 4})
				}


				sprites += 4
			}
		}
	}
}