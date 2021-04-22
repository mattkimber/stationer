package nfo

import (
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

type StationHall struct {
	ID                    int
	BarePlatformSprite    int
	RoofPlatformSprite    int
	RoofBaseSprite        int
	ClassID               string
	ClassName             string
	ObjectName            string
	PlatformConfiguration properties.PlatformLayout
	UseCompanyColour      bool
	MaxLoadState          int
	PlatformHeight        int
	YearAvailable         int
}

func (s *StationHall) SetID(id int) {
	s.ID = id
}

func (s *StationHall) GetID() int {
	return s.ID
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

	innerPlatformSprite := s.GetBaseSpriteNumber() + s.BarePlatformSprite + (direction * 2)
	outerPlatformSprite := s.GetBaseSpriteNumber() + s.BarePlatformSprite + 1 + (direction * 2)

	if supportInner {
		innerPlatformSprite = s.GetBaseSpriteNumber() + s.RoofPlatformSprite + (direction * 2)
	}

	if supportOuter {
		outerPlatformSprite = s.GetBaseSpriteNumber() + s.RoofPlatformSprite + 1 + (direction * 2)
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
		SpriteNumber: TRANSPARENT_SPRITE + s.RoofBaseSprite + (roofSprite * 4) + direction + 2,
	})

	// Add the roof
	result = append(result, properties.BoundingBox{
		XOffset:      0,
		YOffset:      0,
		ZOffset:      14,
		X:            16,
		Y:            16,
		Z:            6,
		SpriteNumber: s.GetBaseSpriteNumber() + s.RoofBaseSprite + (roofSprite * 4) + direction,
	})

	return result
}

func (s *StationHall) WriteToFile(file *output_file.File) {

	if s.MaxLoadState == 0 {
		s.MaxLoadState = DEFAULT_MAX_LOAD_STATE
	}

	if s.PlatformHeight == 0 {
		s.PlatformHeight = DEFAULT_PLATFORM_HEIGHT
	}

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})
	def.AddProperty(&properties.LittleLotsThreshold{Amount: 40})

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
		NumLittleSets: LITTLE_SETS * 2,
		NumLotsSets:   LOTS_SETS,
		SpriteSets:    GetSpriteSets(s.MaxLoadState),
	})

	file.AddElement(&StationSet{
		SetID:         1,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    GetSpriteSets(0),
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
