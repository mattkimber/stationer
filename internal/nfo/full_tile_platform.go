package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
	"github.com/mattkimber/stationer/internal/nfo/sprites"
)

type FullTilePlatform struct {
	ID                   int
	SpriteFilename       string
	ClassID              string
	ClassName            string
	ObjectName           string
	Width                int
	Height               int
	UseCompanyColour     bool
	HasCustomFoundations bool
	YearAvailable        int
}


func (s *FullTilePlatform) SetID(id int) {
	s.ID = id
}

func (s *FullTilePlatform) GetID() int {
	return s.ID
}

func (s *FullTilePlatform) GetBaseSpriteNumber() int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func (s *FullTilePlatform) GetObjects(direction int, idx int) []properties.BoundingBox {
	x, y := 16, 16
	result := make([]properties.BoundingBox, 0)
	base := 0
	if direction == NORTH_SOUTH {
		base = 1
	}
	result = append(result, properties.BoundingBox{X: x, Y: y, Z: 3, SpriteNumber: s.GetBaseSpriteNumber() + base + (idx * 2)})

	return result
}

func (s *FullTilePlatform) WriteToFile(file *output_file.File) {
	s.addSprites(file)

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})

	layoutEntries := make([]properties.LayoutEntry, 0)

	// Add the layouts
	for i := 0; i < 2*16; i++ {
		entry := s.getLayoutEntry(i)
		layoutEntries = append(layoutEntries, entry)
	}

	def.AddProperty(&properties.SpriteLayout{
		Entries: layoutEntries,
	})

	// No pylons or wires
	def.AddProperty(&properties.PylonPlacement{Bitmask: 0})
	def.AddProperty(&properties.WirePlacement{Bitmask: 255})

	// Prevent train entering
	def.AddProperty(&properties.PreventTrainEntryFlag{})

	// Add flag for sprite layout callback
	def.AddProperty(&properties.CallbackFlag{SpriteLayout: true, Availability: s.YearAvailable != 0})

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	yearCallbackID, spriteset := 0, 0

	if s.YearAvailable != 0 {
		yearCallbackID = 20
		file.AddElement(&callbacks.AvailabilityYearCallback{
			SetID:      yearCallbackID,
			HasDecider: false,
			Year:       s.YearAvailable,
		})
	}

	// Add the callback
	file.AddElement(&callbacks.PlatformFenceCallback{YearCallbackID: yearCallbackID})

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

func (s *FullTilePlatform) addSprites(file *output_file.File) {
	// 3 sprites: N, S and both - 3 for Both fences - 7 each for N/S fence combinations
	// 2 directions so all are doubled
	numSprites := 2 * 16
	file.AddElement(&sprites.Spritesets{ID: 0, NumSets: 1, NumSprites: numSprites})

	filename := fmt.Sprintf("%s_8bpp.png", s.SpriteFilename)

	// Non-fence sprites
	file.AddElement(&sprites.Sprites{
		GetBufferStopSprite(filename, 0, false),
		GetBufferStopSprite(filename, 1, false),
	})

	// Do fences
	fenceElements := []string{"a", "b", "ab", "d", "ad", "bd", "abd", "c", "ac", "bc", "abc", "cd", "acd", "bcd", "abcd"}

	for _, fenceElement := range fenceElements {
		filename := fmt.Sprintf("%s_fence_%s_8bpp.png", s.SpriteFilename, fenceElement)

		file.AddElement(&sprites.Sprites{
			GetBufferStopSprite(filename, 0, false),
			GetBufferStopSprite(filename, 1, false),
		})
	}

}

func (s *FullTilePlatform) getLayoutEntry(idx int) properties.LayoutEntry {
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
