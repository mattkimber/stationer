package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

type Station struct {
	ID int
	SpriteFilename string
	ClassID string
	ClassName string
	ObjectName string
	PlatformConfiguration properties.PlatformLayout
	UseCompanyColour bool
}

const (
	SPRITE_WIDTH_WITH_PADDING = 72
	SPRITE_WIDTH = 64
	SPRITE_HEIGHT = 55

	CUSTOM_SPRITE = 0x42D
	COMPANY_COLOUR_SPRITE = 0x842D
	TRANSPARENT_SPRITE = 0x322442D

	MAX_LOAD_STATE = 6
	LITTLE_SETS = 4
	LOTS_SETS = 3
)

func GetSprite(filename string, num int) Sprite {
	xrel := -(SPRITE_WIDTH/2)-10

	if num >= 2 {
		xrel = 11-(SPRITE_WIDTH/2)
	}

	return Sprite{
		Filename: filename,
		X:        SPRITE_WIDTH_WITH_PADDING*num,
		Y:        0,
		Width:    SPRITE_WIDTH,
		Height:   SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     -(SPRITE_HEIGHT/2)-1,
	}
}

func (s *Station) GetBaseSpriteNumber() int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func (s *Station) WriteToFile(file *File) {

	file.AddElement(&Spritesets{NumSets: MAX_LOAD_STATE + 1, NumSprites: 4})

	for i := 0; i <= MAX_LOAD_STATE; i++ {
		filename := fmt.Sprintf("%s_%d_8bpp.png", s.SpriteFilename, i)

		file.AddElement(&Sprites{
			GetSprite(filename, 0),
			GetSprite(filename, 1),
			GetSprite(filename, 2),
			GetSprite(filename, 3),
		})
	}

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})
	def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	def.AddProperty(&properties.SpriteLayout{
		EastWest:   properties.SpriteDirection{
			GroundSprite: 1012,
			Sprites: []properties.BoundingBox{
				{YOffset: 16 - 5, X: 16, Y: 5, Z: 2, SpriteNumber: s.GetBaseSpriteNumber() + 0},
				{X: 16, Y: 5, Z: 3, SpriteNumber: s.GetBaseSpriteNumber() + 1},
			},
		},
		NorthSouth: properties.SpriteDirection{
			GroundSprite: 1011,
			Sprites: []properties.BoundingBox{
				{XOffset: 16 - 5, X: 5, Y: 16, Z: 3, SpriteNumber: s.GetBaseSpriteNumber() + 2},
				{X: 5, Y: 16, Z: 3, SpriteNumber: s.GetBaseSpriteNumber() + 3},
			},
		},
	})

	defaultLayout := properties.PlatformLayout{}
	if s.PlatformConfiguration != defaultLayout {
		def.AddProperty(&properties.AllowedLengths{Bitmask: s.PlatformConfiguration.Lengths})
		def.AddProperty(&properties.AllowedPlatforms{Bitmask: s.PlatformConfiguration.Platforms})
	}

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: LITTLE_SETS*3,
		NumLotsSets:   LOTS_SETS*3,
		SpriteSets:    []int{0,1,1,1,2,2,2,2,3,3,3,3,3,4,4,4,5,5,5,5,6},
	})

	file.AddElement(&StationSet{
		SetID:         1,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	file.AddElement(&GraphicSetAssignment{
		IDs:               []int {s.ID},
		CargoSpecificSets: []CargoToSet{{
			CargoType: 0,
			Set:       0,
		}},
		DefaultSet:        1,
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
		Text:          s.ClassName,
	})
}