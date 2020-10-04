package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

type Building struct {
	ID                    int
	SpriteFilename        string
	ClassID               string
	ClassName             string
	ObjectName            string
	Width int
	Height int
	UseCompanyColour      bool
	HasCustomFoundations  bool
}

const (
	BUILDING_SPRITE_WIDTH_WITH_PADDING = 72
	BUILDING_SPRITE_WIDTH              = 64
	BUILDING_SPRITE_HEIGHT             = 78
)

func GetBuildingSprite(filename string, num int, swap bool) Sprite {
	xrel := 1-(BUILDING_SPRITE_WIDTH / 2)
	yrel := -(BUILDING_SPRITE_HEIGHT / 2) - 6

	if swap {
		xrel = 11 - (BUILDING_SPRITE_WIDTH / 2)
	}

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

func (s *Building) GetBaseSpriteNumber() int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func (s *Building) GetObjects(direction int) []properties.BoundingBox {
	x, y := 16, 16
	result := make([]properties.BoundingBox, 0)
	base := 0
	if direction == NORTH_SOUTH {
		base = 1
	}
	result = append(result, properties.BoundingBox{X: x, Y: y, Z: 16, SpriteNumber: s.GetBaseSpriteNumber() + base})

	return result
}

func (s *Building) WriteToFile(file *File) {

	platformSprites := 2

	file.AddElement(&Spritesets{ID: 0, NumSets: 1, NumSprites: platformSprites})

	filename := fmt.Sprintf("%s_8bpp.png", s.SpriteFilename)

	// Non-fence sprites
	file.AddElement(&Sprites{
		GetBuildingSprite(filename, 0, false),
		GetBuildingSprite(filename, 1, false),
	})

	// Foundation sprites
	if s.HasCustomFoundations {
		file.AddElement(&Spritesets{ID: 1, NumSets: 1, NumSprites: 2})

		filename := fmt.Sprintf("%s_foundation_8bpp.png", s.SpriteFilename)

		// Non-fence sprites
		file.AddElement(&Sprites{
			GetBuildingSprite(filename, 0, false),
			GetBuildingSprite(filename, 1, false),
		})
	}

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})

	// This is irrelevant?
	// def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	layoutEntries := make([]properties.LayoutEntry, 0)

	entry := properties.LayoutEntry{
		EastWest: properties.SpriteDirection{
			GroundSprite: 3981,
			Sprites:      s.GetObjects(EAST_WEST),
		},
		NorthSouth: properties.SpriteDirection{
			GroundSprite: 3981,
			Sprites:      s.GetObjects(NORTH_SOUTH),
		},
	}

	layoutEntries = append(layoutEntries, entry)

	def.AddProperty(&properties.SpriteLayout{
		Entries: layoutEntries,
	})

	// Limit to 1x1 layout
	def.AddProperty(&properties.AllowedLengths{Bitmask: properties.PlatformBitmask{Enable1: true}})
	def.AddProperty(&properties.AllowedPlatforms{Bitmask: properties.PlatformBitmask{Enable1: true}})

	// Prevent train entering
	def.AddProperty(&properties.PreventTrainEntryFlag{})

	// Add flags
	def.AddProperty(&properties.GeneralFlag{HasCustomFoundations: s.HasCustomFoundations})

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	spriteset := 0
	if s.HasCustomFoundations {
		file.AddElement(&StationSet{
			SetID:         1,
			NumLittleSets: 0,
			NumLotsSets:   1,
			SpriteSets:    []int{1},
		})

		file.AddElement(&callbacks.FoundationCallback{
			SetID:            2,
		})

		spriteset = 2
	}


	file.AddElement(&GraphicSetAssignment{
		IDs: []int{s.ID},
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
