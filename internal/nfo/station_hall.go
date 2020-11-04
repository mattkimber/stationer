package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

type StationHall struct {
	ID                    int
	SpriteFilename        string
	ClassID               string
	ClassName             string
	ObjectName            string
	PlatformConfiguration properties.PlatformLayout
	UseCompanyColour      bool
	MaxLoadState          int
	PlatformHeight        int
	RoofType              string
	YearAvailable         int
}

func GetRoofSprite(filename string, num int) Sprite {
	xrel := 1 - (BUILDING_SPRITE_WIDTH / 2)
	yrel := -(BUILDING_SPRITE_HEIGHT / 2) - 6

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

func (s *StationHall) GetBaseSpriteNumber() int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func (s *StationHall) GetObjects(direction int, supportOuter bool, supportInner bool, roofSprite int) []properties.BoundingBox {
	yOffset := 16 - 5
	xOffset := 0
	x, y := 16, 5
	base := 0

	if direction == NORTH_SOUTH {
		xOffset = 16 - 5
		yOffset = 0
		x = 5
		y = 16
		base = base + 2
	}

	result := make([]properties.BoundingBox, 0)

	innerPlatformSprite := s.GetBaseSpriteNumber() + base + 0
	outerPlatformSprite := s.GetBaseSpriteNumber() + base + 1

	if supportInner {
		innerPlatformSprite += 4
	}

	if supportOuter {
		outerPlatformSprite += 4
	}

	// Add the base tiles
	result = append(result, properties.BoundingBox{YOffset: yOffset, XOffset: xOffset, X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: innerPlatformSprite})
	result = append(result, properties.BoundingBox{X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: outerPlatformSprite})

	// Add the glass
	result = append(result, properties.BoundingBox{
		XOffset:      0,
		YOffset:      0,
		ZOffset:      14,
		X:            16,
		Y:            16,
		Z:            6,
		SpriteNumber: TRANSPARENT_SPRITE + 8 + (roofSprite * 4) + direction + 2,
	})

	// Add the roof
	result = append(result, properties.BoundingBox{
		XOffset:      0,
		YOffset:      0,
		ZOffset:      14,
		X:            16,
		Y:            16,
		Z:            6,
		SpriteNumber: s.GetBaseSpriteNumber() + 8 + (roofSprite * 4) + direction,
	})

	return result
}

func (s *StationHall) WriteToFile(file *File) {

	if s.MaxLoadState == 0 {
		s.MaxLoadState = DEFAULT_MAX_LOAD_STATE
	}

	if s.PlatformHeight == 0 {
		s.PlatformHeight = DEFAULT_PLATFORM_HEIGHT
	}

	platformSprites := 8

	file.AddElement(&Spritesets{ID: 0, NumSets: s.MaxLoadState + 1, NumSprites: platformSprites + (3 * 4)})

	for i := 0; i <= s.MaxLoadState; i++ {
		// Sprites without roof support
		filename := fmt.Sprintf("%s_%d_8bpp.png", s.SpriteFilename, i)

		file.AddElement(&Sprites{
			GetSprite(filename, 0, false),
			GetSprite(filename, 1, false),
			GetSprite(filename, 2, true),
			GetSprite(filename, 3, true),
		})

		// Sprites with roof support
		filename = fmt.Sprintf("%s_%d_roof_8bpp.png", s.SpriteFilename, i)

		file.AddElement(&Sprites{
			GetSprite(filename, 0, false),
			GetSprite(filename, 1, false),
			GetSprite(filename, 2, true),
			GetSprite(filename, 3, true),
		})

		roofSubtypes := []string{"single", "double_n", "double_s"}

		for _, subtype := range roofSubtypes {
			// 3 roof sprites (single, N, S) - repeated as we have multiple load states
			filename = fmt.Sprintf("roof_%s_%s_8bpp.png", s.RoofType, subtype)
			file.AddElement(&Sprites{
				GetRoofSprite(filename, 0),
				GetRoofSprite(filename, 1),
			})

			// And again for the glass
			filename = fmt.Sprintf("roof_%s_%s_glass_8bpp.png", s.RoofType, subtype)
			file.AddElement(&Sprites{
				GetRoofSprite(filename, 0),
				GetRoofSprite(filename, 1),
			})
		}

	}

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})
	def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	if s.YearAvailable != 0 {
		def.AddProperty(&properties.CallbackFlag{Availability: s.YearAvailable != 0})
	}

	// Default layouts as per original OpenTTD stations
	layoutEntries := []properties.LayoutEntry{
		// Single platform
		{
			EastWest: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_EW,
				Sprites:      s.GetObjects(EAST_WEST, true, true, 0),
			},
			NorthSouth: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_NS,
				Sprites:      s.GetObjects(NORTH_SOUTH, true, true, 0),
			},
		},
		// Single platform with building (here is identical to single platform)
		{
			EastWest: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_EW,
				Sprites:      s.GetObjects(EAST_WEST, true, true, 0),
			},
			NorthSouth: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_NS,
				Sprites:      s.GetObjects(NORTH_SOUTH, true, true, 0),
			},
		},
		// Roof, L
		{
			EastWest: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_EW,
				Sprites:      s.GetObjects(EAST_WEST, true, false, 1),
			},
			NorthSouth: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_NS,
				Sprites:      s.GetObjects(NORTH_SOUTH, true, false, 1),
			},
		},
		// Roof, R
		{
			EastWest: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_EW,
				Sprites:      s.GetObjects(EAST_WEST, false, true, 2),
			},
			NorthSouth: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_NS,
				Sprites:      s.GetObjects(NORTH_SOUTH, false, true, 2),
			},
		},
	}

	def.AddProperty(&properties.SpriteLayout{
		Entries: layoutEntries,
	})

	defaultLayout := properties.PlatformLayout{}
	if s.PlatformConfiguration != defaultLayout {
		def.AddProperty(&properties.AllowedLengths{Bitmask: s.PlatformConfiguration.Lengths})
		def.AddProperty(&properties.AllowedPlatforms{Bitmask: s.PlatformConfiguration.Platforms})
	}

	passengerCargoSet := 0
	otherCargoSet := 1

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: LITTLE_SETS * 3,
		NumLotsSets:   LOTS_SETS * 3,
		SpriteSets:    GetSpriteSets(s.MaxLoadState),
	})

	file.AddElement(&StationSet{
		SetID:         1,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	// The callback definition has to come *after* the set definitions or it will be referencing sets from the
	// previous station item.
	setID := 0
	if s.YearAvailable != 0 {
		setID = 8
		file.AddElement(&callbacks.AvailabilityYearCallback{
			SetID:            setID,
			HasDecider:       true,
			Year:             s.YearAvailable,
			DefaultSpriteset: 1,
		})

		otherCargoSet = setID
	}

	file.AddElement(&GraphicSetAssignment{
		IDs: []int{s.ID},
		CargoSpecificSets: []CargoToSet{{
			CargoType: 0,
			Set:       passengerCargoSet,
		}},
		DefaultSet: otherCargoSet,
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
