package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

type BufferStop struct {
	ID                   int
	SpriteFilename       string
	ClassID              string
	ClassName            string
	ObjectName           string
	Width                int
	Height               int
	UseCompanyColour     bool
	HasCustomFoundations bool
}

const (
	BUFFERSTOP_SPRITE_WIDTH_WITH_PADDING = 72
	BUFFERSTOP_SPRITE_WIDTH              = 64
	BUFFERSTOP_SPRITE_HEIGHT             = 55
)

func (s *BufferStop) GetBaseSpriteNumber() int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func GetBufferStopSprite(filename string, num int, swap bool) Sprite {
	xrel := 1 - (BUFFERSTOP_SPRITE_WIDTH / 2)
	yrel := 5 - (BUFFERSTOP_SPRITE_HEIGHT / 2)

	if swap {
		xrel = 11 - (BUFFERSTOP_SPRITE_WIDTH / 2)
	}

	return Sprite{
		Filename: filename,
		X:        BUFFERSTOP_SPRITE_WIDTH_WITH_PADDING * num,
		Y:        0,
		Width:    BUFFERSTOP_SPRITE_WIDTH,
		Height:   BUFFERSTOP_SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     yrel,
	}
}

func (s *BufferStop) GetObjects(direction int, idx int) []properties.BoundingBox {
	x, y := 16, 16
	result := make([]properties.BoundingBox, 0)
	base := 0
	if direction == NORTH_SOUTH {
		base = 1
	}
	result = append(result, properties.BoundingBox{X: x, Y: y, Z: 3, SpriteNumber: s.GetBaseSpriteNumber() + base + (idx * 2)})

	return result
}

func (s *BufferStop) WriteToFile(file *File) {
	s.addSprites(file)

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})

	// This is irrelevant?
	// def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	layoutEntries := make([]properties.LayoutEntry, 0)

	// Add the layouts
	for i := 0; i < (3 + 3 + (7 * 2)); i++ {
		entry := s.getLayoutEntry(i)
		layoutEntries = append(layoutEntries, entry)
	}

	def.AddProperty(&properties.SpriteLayout{
		Entries: layoutEntries,
	})

	def.AddProperty(&properties.AllowedLengths{Bitmask: properties.PlatformBitmask{Enable1: true}})

	// No pylons or wires
	def.AddProperty(&properties.PylonPlacement{Bitmask: 0})
	def.AddProperty(&properties.WirePlacement{Bitmask: 255})

	// Prevent train entering
	def.AddProperty(&properties.PreventTrainEntryFlag{})

	// Add flag for sprite layout callback
	def.AddProperty(&properties.CallbackFlag{SpriteLayout: true})

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	spriteset := 0

	// Add the callback
	file.AddElement(&callbacks.BufferStopCallback{})

	file.AddElement(&GraphicSetAssignment{
		IDs:        []int{s.ID},
		DefaultSet: spriteset,
	})

	file.AddElement(&TextString{
		LanguageFile:   255,
		StationId:      s.ID,
		TextStringType: TextStringTypeStationName,
		Text:           s.ObjectName,
	})

	file.AddElement(&TextString{
		LanguageFile:   255,
		StationId:      s.ID,
		TextStringType: TextStringTypeClassName,
		Text:           s.ClassName,
	})
}

func (s *BufferStop) addSprites(file *File) {
	// 3 sprites: N, S and both - 3 for Both fences - 7 each for N/S fence combinations
	// 2 directions so all are doubled
	numSprites := 2 * (3 + 3 + (7 * 2))
	file.AddElement(&Spritesets{ID: 0, NumSets: 1, NumSprites: numSprites})

	elements := []string{"s", "n", "both"}
	for _, element := range elements {

		filename := fmt.Sprintf("%s_%s_8bpp.png", s.SpriteFilename, element)

		// Non-fence sprites
		file.AddElement(&Sprites{
			GetBufferStopSprite(filename, 0, false),
			GetBufferStopSprite(filename, 1, false),
		})

		// Do fences
		fenceElements := []string{"a", "b", "ab"}

		if element == "s" {
			fenceElements = []string{"a", "b", "ab", "d", "ad", "bd", "abd"}
		} else if element == "n" {
			fenceElements = []string{"a", "b", "ab", "c", "ac", "bc", "abc"}
		}

		for _, fenceElement := range fenceElements {
			filename := fmt.Sprintf("%s_%s_fence_%s_8bpp.png", s.SpriteFilename, element, fenceElement)

			file.AddElement(&Sprites{
				GetBufferStopSprite(filename, 0, false),
				GetBufferStopSprite(filename, 1, false),
			})
		}
	}
}

func (s *BufferStop) getLayoutEntry(idx int) properties.LayoutEntry {
	entry := properties.LayoutEntry{
		EastWest: properties.SpriteDirection{
			GroundSprite: GROUND_SPRITE_RAIL_EW,
			Sprites:      s.GetObjects(EAST_WEST, idx),
		},
		NorthSouth: properties.SpriteDirection{
			GroundSprite: GROUND_SPRITE_RAIL_NS,
			Sprites:      s.GetObjects(NORTH_SOUTH, idx),
		},
	}
	return entry
}
