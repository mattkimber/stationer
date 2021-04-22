package sprites

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
)

type StationSprite struct {
	Filename            string
	HasFences           bool
	MaxLoadState        int
	DedicatedFlipSprite bool // If the sprite has its own file for flipped sprites
	SingleSided         bool // Only one side is needed - the rear side is never displayed
	IsStatic            bool // Is a single static sprite for all load states
	IsRailFence         bool // Use the class filename even if the sprite is static
}

type StationSprites struct {
	BaseFilename     string
	Sprites          []StationSprite
	SpriteMap        map[string]int // Will be populated after SetStatistics is called
	LastSpriteNumber int            // Will be populated after SetStatistics is called
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

func (s *StationSprites) SetStatistics() {
	total := 0
	s.SpriteMap = make(map[string]int)

	for _, sprite := range s.Sprites {
		s.SpriteMap[sprite.Filename] = total

		total += 4

		if sprite.HasFences || sprite.IsRailFence {
			// Add another 4 for the fences
			total += 4
		}

	}

	s.LastSpriteNumber = total
}

func (s *StationSprites) WriteToFile(file *output_file.File, loadState int) {

	for _, spr := range s.Sprites {
		// Populate the map of where sprites begin and end

		filename := fmt.Sprintf("%s_%s_%d_8bpp.png", s.BaseFilename, spr.Filename, loadState)
		if spr.Filename == "roof" {
			// Hack to override the different naming scheme of roofs
			filename = fmt.Sprintf("%s_empty_%d_%s_8bpp.png", s.BaseFilename, loadState, spr.Filename)
		}

		if spr.IsStatic {
			// Ignore all load states and platform types
			// This is slightly wasteful as we repeat the same sprite multiple times.
			if spr.IsRailFence {
				// Note the fences are the other way round from normal sprite naming convention
				filename = fmt.Sprintf("%s_%s_8bpp.png", spr.Filename, s.BaseFilename)
			} else {
				filename = fmt.Sprintf("%s_8bpp.png", spr.Filename)
			}
		}

		filenameFlip := filename
		if spr.DedicatedFlipSprite {
			filenameFlip = fmt.Sprintf("%s_%s_flip_%d_8bpp.png", s.BaseFilename, spr.Filename, loadState)
		}

		if loadState <= spr.MaxLoadState {
			if spr.IsRailFence {
				// Add 4x transparent sprites for the "non-fence" appearance of rail-side
				// fences on single side platforms
				file.AddElement(&Sprites{Sprite{}, Sprite{}, Sprite{}, Sprite{}})
			}

			if spr.SingleSided {
				// Non-fence sprites for single-sided object
				// Uses blank sprites in between to keep the same relative sprite offsets
				file.AddElement(&Blank{Size: 1})
				file.AddElement(&Sprites{s.GetSprite(filename, 1, false)})
				file.AddElement(&Blank{Size: 1})
				file.AddElement(&Sprites{s.GetSprite(filenameFlip, 3, true)})
			} else {
				// Non-fence sprites
				file.AddElement(&Sprites{
					s.GetSprite(filename, 0, false),
					s.GetSprite(filename, 1, false),
					s.GetSprite(filenameFlip, 2, true),
					s.GetSprite(filenameFlip, 3, true),
				})
			}

		} else {
			// Add blank pseudosprites
			file.AddElement(&Blank{Size: 4, Name: filename})
			if spr.IsRailFence {
				// Blank equivalent of transparent fence sprites
				file.AddElement(&Blank{Size: 4, Name: filename})
			}
		}

		// Fence sprites
		if spr.HasFences {
			fenceFilename := fmt.Sprintf("%s_%s_%d_fence_8bpp.png", s.BaseFilename, spr.Filename, loadState)
			if loadState <= spr.MaxLoadState {
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
		}
	}
}
