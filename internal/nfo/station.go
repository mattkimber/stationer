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
	MaxLoadState int
	AdditionalObjects []AdditionalObject
	PlatformHeight int
}

const (
	SPRITE_WIDTH_WITH_PADDING = 72
	SPRITE_WIDTH = 64
	SPRITE_HEIGHT = 55
	DEFAULT_PLATFORM_HEIGHT = 3

	CUSTOM_SPRITE = 0x42D
	COMPANY_COLOUR_SPRITE = 0x842D
	TRANSPARENT_SPRITE = 0x322442D

	DEFAULT_MAX_LOAD_STATE = 6
	LITTLE_SETS = 4
	LOTS_SETS = 3

	EAST_WEST = 0
	NORTH_SOUTH = 1
)

func GetSprite(filename string, num int, swap bool) Sprite {
	xrel := -(SPRITE_WIDTH/2)-10

	if swap {
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

func GetSpriteSets(max int) []int {
	switch max {
	case 5:
		return []int{0,1,1,1,2,2,2,2,3,3,3,3,3,4,4,4,5,5,5,5,5}
	case 6:
		return []int{0,1,1,1,2,2,2,2,3,3,3,3,3,4,4,4,5,5,5,5,6}
	}

	return []int{0}
}

func (s *Station) GetObjects(direction int) []properties.BoundingBox {
	yOffset := 16 - 5
	xOffset := 0
	x, y := 16, 5
	base := 0

	if direction == NORTH_SOUTH {
		xOffset = 16 - 5
		yOffset = 0
		x = 5
		y = 16
		base = 2
	}

	result := []properties.BoundingBox{
		{YOffset: yOffset, XOffset: xOffset, X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: s.GetBaseSpriteNumber() + base + 0},
		{X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: s.GetBaseSpriteNumber() + base + 1},
	}

	for idx, obj := range s.AdditionalObjects {
		x, y := obj.SizeX, obj.SizeY
		xOffset, yOffset := obj.X, obj.Y

		if direction == NORTH_SOUTH {
			x, y = obj.SizeY, obj.SizeX
			xOffset, yOffset = obj.Y, obj.X

		}

		result = append(result, properties.BoundingBox{
			XOffset:      xOffset,
			YOffset:      yOffset,
			ZOffset:      obj.Z,
			X:            x,
			Y:            y,
			Z:            obj.SizeZ,
			SpriteNumber: s.GetBaseSpriteNumber() + 4 + (idx*2) + direction,
		})
	}

	return result
}

func (s *Station) WriteToFile(file *File) {

	if s.MaxLoadState == 0 {
		s.MaxLoadState = DEFAULT_MAX_LOAD_STATE
	}

	if s.PlatformHeight == 0 {
		s.PlatformHeight = DEFAULT_PLATFORM_HEIGHT
	}


	file.AddElement(&Spritesets{NumSets: s.MaxLoadState + 1, NumSprites: 4 + (len(s.AdditionalObjects) * 2)})

	for i := 0; i <= s.MaxLoadState; i++ {
		filename := fmt.Sprintf("%s_%d_8bpp.png", s.SpriteFilename, i)

		file.AddElement(&Sprites{
			GetSprite(filename, 0, false),
			GetSprite(filename, 1, false),
			GetSprite(filename, 2, true),
			GetSprite(filename, 3, true),
		})

		for _, obj := range s.AdditionalObjects {
			filename := fmt.Sprintf("%s_8bpp.png", obj.SpriteFilename)
			file.AddElement(&Sprites{
				GetSprite(filename, 0, obj.InvertDirection != true),
				GetSprite(filename, 1, obj.InvertDirection != false),
			})
		}
	}


	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})
	def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	def.AddProperty(&properties.SpriteLayout{
		EastWest:   properties.SpriteDirection{
			GroundSprite: 1012,
			Sprites: s.GetObjects(EAST_WEST),
		},
		NorthSouth: properties.SpriteDirection{
			GroundSprite: 1011,
			Sprites: s.GetObjects(NORTH_SOUTH),
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
		SpriteSets:    GetSpriteSets(s.MaxLoadState),
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