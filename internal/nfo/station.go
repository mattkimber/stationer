package nfo

import (
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

type Station struct {
	ID                    int
	BaseSpriteID          int
	RandomSpriteIDs       []int
	ClassID               string
	ClassName             string
	ObjectName            string
	PlatformConfiguration properties.PlatformLayout
	UseCompanyColour      bool
	MaxLoadState          int
	AdditionalObjects     []AdditionalObject
	PlatformHeight        int
	HasFences             bool
	InnerPlatform         bool
	OuterPlatform         bool
	YearAvailable         int
	OverrideOuter         bool // Whether to use an override sprite for the outer platforms
	OuterPlatformSprite   int  // and if so, which sprite.
	HasLargeCentralObject bool // If the central object covers two tiles
	ObjectIsSingleSided   bool // if the central object is "single sided" - only present on one platform tile.
}

const (
	DEFAULT_PLATFORM_HEIGHT = 3

	CUSTOM_SPRITE         = 0x42D
	COMPANY_COLOUR_SPRITE = 0x842D
	TRANSPARENT_SPRITE    = 0x322442D

	DEFAULT_MAX_LOAD_STATE = 6
	LITTLE_SETS            = 4
	LOTS_SETS              = 3

	EAST_WEST   = 0
	NORTH_SOUTH = 1

	GROUND_SPRITE_RAIL_EW = 1012
	GROUND_SPRITE_RAIL_NS = 1011
)

func (s *Station) GetBaseSpriteNumber() int {
	// The sprite number is a relative offset from the spriteset.
	// e.g. even if you are using set ID 20, the base sprite in it
	// still has number 0.
	//
	// However in this set we put all of the station graphics in
	// a single spriteset so we can mix and match between them
	// in our layouts - hence there is still an offset here, which
	// is the offset *within* the spriteset.
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE + s.BaseSpriteID
	}

	return CUSTOM_SPRITE + s.BaseSpriteID
}

func (s *Station) GetRandomSpriteNumber(number int) int {

	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE + s.RandomSpriteIDs[number]
	}

	return CUSTOM_SPRITE + s.RandomSpriteIDs[number]
}

func (s *Station) GetOuterPlatformSpriteNumber() int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE + s.OuterPlatformSprite
	}

	return CUSTOM_SPRITE + s.OuterPlatformSprite
}

func GetSpriteSets(max int) []int {
	switch max {
	case 5:
		return []int{0, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 5}
	case 6:
		return []int{0, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 6}
	}

	return []int{0}
}

func (s *Station) GetObjects(direction int, fenceInside, fenceOutside bool, iteration int) []properties.BoundingBox {
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

	platform := iteration % 3
	randomChoice := iteration / 3

	if s.OuterPlatform || (platform >= 1 && s.HasLargeCentralObject) {
		baseSprite := s.GetBaseSpriteNumber()

		if len(s.RandomSpriteIDs) > 0 {
			baseSprite = s.GetRandomSpriteNumber(randomChoice)
		}

		if s.OverrideOuter && (platform == 0 || s.ObjectIsSingleSided) {
			baseSprite = s.GetOuterPlatformSpriteNumber()
		}

		// We only do the fence offsets if we are on the "outside" tile, not the central object tile.
		// This is so when players chop stations up into "unnatural" tiles (e.g. removing one side of a
		// 2-platform object) OpenTTD doesn't get confused looking for an "object+fence" sprite that doesn't exist.
		if fenceOutside && platform == 0 {
			result = append(result, properties.BoundingBox{YOffset: yOffset, XOffset: xOffset, X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: baseSprite + base + 4})
		} else {
			result = append(result, properties.BoundingBox{YOffset: yOffset, XOffset: xOffset, X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: baseSprite + base + 0})
		}
	}

	if s.InnerPlatform && !(!s.OuterPlatform && s.HasLargeCentralObject && platform == 1) {
		baseSprite := s.GetBaseSpriteNumber()

		if len(s.RandomSpriteIDs) > 0 {
			baseSprite = s.GetRandomSpriteNumber(randomChoice)
		}

		if s.OverrideOuter && platform == 1 {
			baseSprite = s.GetOuterPlatformSpriteNumber()
		}

		if fenceInside && (platform == 1 || !s.HasLargeCentralObject) {
			result = append(result, properties.BoundingBox{X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: baseSprite + base + 5})
		} else {
			result = append(result, properties.BoundingBox{X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: baseSprite + base + 1})
		}
	}

	for _, obj := range s.AdditionalObjects {
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
			SpriteNumber: obj.GetBaseSpriteNumber(s) + direction,
		})
	}

	return result
}

func (s *Station) WriteToFile(file *output_file.File) {

	if s.MaxLoadState == 0 {
		s.MaxLoadState = DEFAULT_MAX_LOAD_STATE
	}

	if s.PlatformHeight == 0 {
		s.PlatformHeight = DEFAULT_PLATFORM_HEIGHT
	}

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})

	if s.HasFences || s.YearAvailable != 0 {
		def.AddProperty(&properties.CallbackFlag{SpriteLayout: s.HasFences, Availability: s.YearAvailable != 0})
	}

	def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	iterations := 1
	if len(s.RandomSpriteIDs) > 0 {
		iterations = len(s.RandomSpriteIDs)
	}

	layoutEntries := make([]properties.LayoutEntry, 0)

	for i := 0; i < iterations; i++ {
		layoutEntries = append(layoutEntries, s.getLayoutEntryForPlatform(i*3)...)
		if s.HasLargeCentralObject {
			// 1 = the "other side" tile
			layoutEntries = append(layoutEntries, s.getLayoutEntryForPlatform(1+(i*3))...)

			// 2 = the "both sides" tile
			layoutEntries = append(layoutEntries, s.getLayoutEntryForPlatform(2+(i*3))...)
		}
	}

	def.AddProperty(&properties.SpriteLayout{
		Entries: layoutEntries,
	})

	defaultLayout := properties.PlatformLayout{}
	if s.PlatformConfiguration != defaultLayout {
		def.AddProperty(&properties.AllowedLengths{Bitmask: s.PlatformConfiguration.Lengths})
		def.AddProperty(&properties.AllowedPlatforms{Bitmask: s.PlatformConfiguration.Platforms})
	}

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

	passengerCargoSet := 0
	otherCargoSet := 1

	yearCallbackID := 0

	if s.YearAvailable != 0 {
		yearCallbackID = 8
		file.AddElement(&callbacks.AvailabilityYearCallback{
			SetID:      yearCallbackID,
			HasDecider: !s.HasFences,
			Year:       s.YearAvailable,
		})

		passengerCargoSet, otherCargoSet = yearCallbackID, yearCallbackID
	}

	randomPassengerSets, randomOtherSets := make([]int, 0), make([]int, 0)

	for i := 0; i < iterations; i++ {
		if s.HasFences {
			file.AddElement(&callbacks.StationFenceCallback{
				SetID:            10 + (i * 50),
				DefaultSpriteSet: 0,
				YearCallbackID:   yearCallbackID,
				HasDecider:       !s.HasLargeCentralObject,
				BaseLayoutOffset: i * 24,
			})

			file.AddElement(&callbacks.StationFenceCallback{
				SetID:            15 + (i * 50),
				DefaultSpriteSet: 1,
				YearCallbackID:   yearCallbackID,
				HasDecider:       !s.HasLargeCentralObject,
				BaseLayoutOffset: i * 24,
			})

			if s.HasLargeCentralObject {
				file.AddElement(&callbacks.StationFenceCallback{
					SetID:            20 + (i * 50),
					DefaultSpriteSet: 0,
					YearCallbackID:   yearCallbackID,
					BaseLayoutOffset: 8 + (i * 24),
				})

				file.AddElement(&callbacks.StationFenceCallback{
					SetID:            25 + (i * 50),
					DefaultSpriteSet: 1,
					YearCallbackID:   yearCallbackID,
					BaseLayoutOffset: 8 + (i * 24),
				})

				file.AddElement(&callbacks.StationFenceCallback{
					SetID:            30 + (i * 50),
					DefaultSpriteSet: 0,
					YearCallbackID:   yearCallbackID,
					BaseLayoutOffset: 16 + (i * 24),
				})

				file.AddElement(&callbacks.StationFenceCallback{
					SetID:            35 + (i * 50),
					DefaultSpriteSet: 1,
					YearCallbackID:   yearCallbackID,
					BaseLayoutOffset: 16 + (i * 24),
				})

				file.AddElement(&callbacks.LargeCentralObjectCallback{
					SetID:            40 + (i * 50),
					OuterCallbackID:  13 + (i * 50),
					InnerCallbackID:  23 + (i * 50),
					MiddleCallbackID: 33 + (i * 50),
					DefaultSpriteSet: 0,
					YearCallbackID:   yearCallbackID,
					HasDecider:       len(s.RandomSpriteIDs) == 0,
				})

				file.AddElement(&callbacks.LargeCentralObjectCallback{
					SetID:            45 + (i * 50),
					OuterCallbackID:  18 + (i * 50),
					InnerCallbackID:  28 + (i * 50),
					MiddleCallbackID: 38 + (i * 50),
					DefaultSpriteSet: 1, // this is needed to prevent stations showing cargo in the purchase menu
					YearCallbackID:   yearCallbackID,
					HasDecider:       len(s.RandomSpriteIDs) == 0,
				})

				passengerCargoSet, otherCargoSet = 40+(i*50), 45+(i*50)

				// The actual set ID of the LCO callback is offset 2 from the base ID
				randomPassengerSets = append(randomPassengerSets, passengerCargoSet+2)
				randomOtherSets = append(randomOtherSets, passengerCargoSet+2)
			} else {
				passengerCargoSet, otherCargoSet = 10+(i*50), 15+(i*50)
			}
		}
	}

	if len(s.RandomSpriteIDs) > 0 {
		// Add the random set decider
		file.AddElement(&callbacks.RandomChoiceCallback{
			SetID:            (iterations + 1) * 50,
			DefaultSpriteSet: 0,
			YearCallbackID:   yearCallbackID,
			ResultIDs:        randomPassengerSets,
			HasDecider:       true,
		})

		file.AddElement(&callbacks.RandomChoiceCallback{
			SetID:            ((iterations + 1) * 50) + 2,
			DefaultSpriteSet: 1, // this is needed to prevent stations showing cargo in the purchase menu
			YearCallbackID:   yearCallbackID,
			ResultIDs:        randomOtherSets,
			HasDecider:       true,
		})

		// The callback decider is 1 + the set ID of the random choice callback
		passengerCargoSet, otherCargoSet = ((iterations+1)*50)+1, ((iterations+1)*50)+3
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

func (s *Station) getLayoutEntryForPlatform(platform int) []properties.LayoutEntry {
	// Default layouts:
	// - Fence (both sides)
	// - Fence (outside)
	// - Fence (inside)
	// - No fences
	layoutEntries := make([]properties.LayoutEntry, 0)
	for i := 0; i < 4; i++ {
		fenceOutside := i >= 2
		fenceInside := i == 1 || i == 3

		if i > 0 && !s.HasFences {
			break
		}

		// Even though not all fence possibilities will be displayed for large objects,
		// we still add the layouts into the station layout so all of the different options
		// have the same size (even if some of the fence possibilities aren't used)
		entry := properties.LayoutEntry{
			EastWest: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_EW,
				Sprites:      s.GetObjects(EAST_WEST, fenceInside, fenceOutside, platform),
			},
			NorthSouth: properties.SpriteDirection{
				GroundSprite: GROUND_SPRITE_RAIL_NS,
				Sprites:      s.GetObjects(NORTH_SOUTH, fenceInside, fenceOutside, platform),
			},
		}

		layoutEntries = append(layoutEntries, entry)
	}
	return layoutEntries
}
