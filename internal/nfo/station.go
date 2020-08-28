package nfo

import "fmt"

type Station struct {
	ID int
	SpriteFilename string
	ClassID string
	ClassName string
	ObjectName string
}

const (
	SPRITE_WIDTH_WITH_PADDING = 72
	SPRITE_WIDTH = 64
	SPRITE_HEIGHT = 55

	CUSTOM_SPRITE = 0x42D
	COMPANY_COLOUR_SPRITE = 0x842D
	TRANSPARENT_SPRITE = 0x322442D
)

func (s *Station) WriteToFile(file *File) {

	filename := fmt.Sprintf("%s_8bpp.png", s.SpriteFilename)

	file.AddElement(&Sprites{
		{
			Filename: filename,
			X:        0,
			Y:        0,
			Width:    SPRITE_WIDTH,
			Height:   SPRITE_HEIGHT,
			XRel:     -(SPRITE_WIDTH/2)-10,
			YRel:     -(SPRITE_HEIGHT/2)-1,
		},
		{
			Filename: filename,
			X:        SPRITE_WIDTH_WITH_PADDING*1,
			Y:        0,
			Width:    SPRITE_WIDTH,
			Height:   SPRITE_HEIGHT,
			XRel:     -(SPRITE_WIDTH/2)-10,
			YRel:     -(SPRITE_HEIGHT/2)-1,
		},
		{
			Filename: filename,
			X:        SPRITE_WIDTH_WITH_PADDING*2,
			Y:        0,
			Width:    SPRITE_WIDTH,
			Height:   SPRITE_HEIGHT,
			XRel:     11-(SPRITE_WIDTH/2),
			YRel:     -(SPRITE_HEIGHT/2)-1,

		},
		{
			Filename: filename,
			X:        SPRITE_WIDTH_WITH_PADDING*3,
			Y:        0,
			Width:    SPRITE_WIDTH,
			Height:   SPRITE_HEIGHT,
			XRel:     11-(SPRITE_WIDTH/2),
			YRel:     -(SPRITE_HEIGHT/2)-1,
		},
	})

	def := &Definition{StationID: s.ID}
	def.AddProperty(&ClassID{ID: s.ClassID})
	def.AddProperty(&SpriteLayout{
		EastWest:   SpriteDirection{
			GroundSprite: 1012,
			Sprites: []BoundingBox{
				{YOffset: 16 - 5, X: 16, Y: 5, Z: 2, SpriteNumber: CUSTOM_SPRITE + 0},
				{X: 16, Y: 5, Z: 3, SpriteNumber: CUSTOM_SPRITE + 1},
			},
		},
		NorthSouth: SpriteDirection{
			GroundSprite: 1011,
			Sprites: []BoundingBox{
				{XOffset: 16 - 5, X: 5, Y: 16, Z: 3, SpriteNumber: CUSTOM_SPRITE + 2},
				{X: 5, Y: 16, Z: 3, SpriteNumber: CUSTOM_SPRITE + 3},
			},
		},
	})

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
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
			CargoType: 254,
			Set:       1,
		}},
		DefaultSet:        0,
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