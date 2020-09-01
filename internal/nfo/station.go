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

	MAX_LOAD_STATE = 6
	LITTLE_SETS = 4
	LOTS_SETS = 2
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
	def.AddProperty(&ClassID{ID: s.ClassID})
	def.AddProperty(&LittleLotsThreshold{Amount: 200})


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
		NumLittleSets: LITTLE_SETS*3,
		NumLotsSets:   LOTS_SETS*2,
		SpriteSets:    []int{0,1,1,1,2,2,2,2,3,3,3,3,4,4,5,5,5,6},
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