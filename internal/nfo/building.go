package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
	"github.com/mattkimber/stationer/internal/nfo/sprites"
)

type Building struct {
	ID               int
	SpriteFilename   string
	ClassID          string
	ClassName        string
	ObjectName       string
	Width            int
	Height           int
	UseCompanyColour bool
	YearAvailable    int
	Reversed         bool
}

const (
	BUILDING_SPRITE_WIDTH_WITH_PADDING = 72
	BUILDING_SPRITE_WIDTH              = 64
	BUILDING_SPRITE_HEIGHT             = 78
)

func GetBuildingSprite(filename string, num int) sprites.Sprite {
	xrel := 1 - (BUILDING_SPRITE_WIDTH / 2)
	yrel := -(BUILDING_SPRITE_HEIGHT / 2) - 6

	return sprites.Sprite{
		Filename: filename,
		X:        BUILDING_SPRITE_WIDTH_WITH_PADDING * num,
		Y:        0,
		Width:    BUILDING_SPRITE_WIDTH,
		Height:   BUILDING_SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     yrel,
	}
}

func (s *Building) GetBaseSpriteNumber() int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func (s *Building) GetObjects(direction int, idx int) []properties.BoundingBox {
	x, y := 16, 16
	result := make([]properties.BoundingBox, 0)
	base := 0
	if direction == NORTH_SOUTH {
		base = 1
	}
	result = append(result, properties.BoundingBox{X: x, Y: y, Z: 16, SpriteNumber: s.GetBaseSpriteNumber() + base + (idx * 2)})

	return result
}

func (s *Building) WriteToFile(file *output_file.File) {
	// Set default width
	if s.Width == 0 {
		s.Width = 1
	}

	s.addSprites(file)

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})

	// This is irrelevant?
	// def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	layoutEntries := make([]properties.LayoutEntry, 0)

	// Add the layouts
	for i := 0; i < s.Width; i++ {
		entry := s.getLayoutEntry(i)
		layoutEntries = append(layoutEntries, entry)
	}

	def.AddProperty(&properties.SpriteLayout{
		Entries: layoutEntries,
	})

	// Limited to 1xn layout
	def.AddProperty(&properties.AllowedLengths{Bitmask: properties.PlatformBitmask{
		Enable1: s.Width == 1,
		Enable2: s.Width == 2,
	}})
	def.AddProperty(&properties.AllowedPlatforms{Bitmask: properties.PlatformBitmask{Enable1: true}})

	// No pylons or wires
	def.AddProperty(&properties.PylonPlacement{Bitmask: 0})
	def.AddProperty(&properties.WirePlacement{Bitmask: 255})

	// Prevent train entering
	def.AddProperty(&properties.PreventTrainEntryFlag{})

	// If this is a multi-tile station or has an availability year it will need a callback for its sprite layout
	if s.Width > 1 || s.YearAvailable != 0 {
		def.AddProperty(&properties.CallbackFlag{SpriteLayout: s.Width > 1, Availability: s.YearAvailable != 0})
	}

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	spriteset := 0
	yearCallbackID := 0

	if s.YearAvailable != 0 {
		yearCallbackID, spriteset = 10, 10
		file.AddElement(&callbacks.AvailabilityYearCallback{
			SetID:      yearCallbackID,
			HasDecider: s.Width <= 1,
			Year:       s.YearAvailable,
		})
	}

	if s.Width > 1 {
		spriteset = 11

		// Add the callback for building tile selection
		file.AddElement(&callbacks.MultiTileBuildingCallback{
			SetID:          spriteset,
			Length:         s.Width,
			YearCallbackID: yearCallbackID,
		})
	}

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

func (s *Building) addSprites(file *output_file.File) {
	buildingSprites := 2 * s.Width
	spriteOffset := 0
	if s.Reversed {
		spriteOffset = 2
	}

	file.AddElement(&sprites.Spritesets{ID: 0, NumSets: 1, NumSprites: buildingSprites})



	// Non-fence sprites
	if !s.Reversed {
		filename := fmt.Sprintf("%s_8bpp.png", s.SpriteFilename)

		file.AddElement(&sprites.Sprites{
			GetBuildingSprite(filename, spriteOffset+0),
			GetBuildingSprite(filename, spriteOffset+1),
		})

		for i := 2; i <= s.Width; i++ {
			// Additional sprites for long buildings
			filename = fmt.Sprintf("%s_%d_8bpp.png", s.SpriteFilename, i)

			file.AddElement(&sprites.Sprites{
				GetBuildingSprite(filename, spriteOffset+0),
				GetBuildingSprite(filename, spriteOffset+1),
			})
		}
	} else {
		// Reversed sprites go in the opposite order
		for i := s.Width; i >= 2; i-- {
			// Additional sprites for long buildings
			filename := fmt.Sprintf("%s_%d_8bpp.png", s.SpriteFilename, i)

			file.AddElement(&sprites.Sprites{
				GetBuildingSprite(filename, spriteOffset+0),
				GetBuildingSprite(filename, spriteOffset+1),
			})
		}

		filename := fmt.Sprintf("%s_8bpp.png", s.SpriteFilename)

		file.AddElement(&sprites.Sprites{
			GetBuildingSprite(filename, spriteOffset+0),
			GetBuildingSprite(filename, spriteOffset+1),
		})
	}


}

func (s *Building) getLayoutEntry(idx int) properties.LayoutEntry {
	entry := properties.LayoutEntry{
		EastWest: properties.SpriteDirection{
			GroundSprite: 3981,
			// East-West sprites are assembled in the "wrong" order so that
			// multi tile stations are the correct way round when displayed
			Sprites: s.GetObjects(EAST_WEST, s.Width-(idx+1)),
		},
		NorthSouth: properties.SpriteDirection{
			GroundSprite: 3981,
			Sprites:      s.GetObjects(NORTH_SOUTH, idx),
		},
	}
	return entry
}
