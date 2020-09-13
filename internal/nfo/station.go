package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

type Station struct {
	ID                    int
	SpriteFilename        string
	ClassID               string
	ClassName             string
	ObjectName            string
	PlatformConfiguration properties.PlatformLayout
	UseCompanyColour      bool
	MaxLoadState          int
	AdditionalObjects     []AdditionalObject
	PlatformHeight        int
	HasFences             bool
}

const (
	SPRITE_WIDTH_WITH_PADDING = 72
	SPRITE_WIDTH              = 64
	SPRITE_HEIGHT             = 55
	DEFAULT_PLATFORM_HEIGHT   = 3

	CUSTOM_SPRITE         = 0x42D
	COMPANY_COLOUR_SPRITE = 0x842D
	TRANSPARENT_SPRITE    = 0x322442D

	DEFAULT_MAX_LOAD_STATE = 6
	LITTLE_SETS            = 4
	LOTS_SETS              = 3

	EAST_WEST   = 0
	NORTH_SOUTH = 1
)

func GetSprite(filename string, num int, swap bool) Sprite {
	xrel := -(SPRITE_WIDTH / 2) - 10

	if swap {
		xrel = 11 - (SPRITE_WIDTH / 2)
	}

	return Sprite{
		Filename: filename,
		X:        SPRITE_WIDTH_WITH_PADDING * num,
		Y:        0,
		Width:    SPRITE_WIDTH,
		Height:   SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     -(SPRITE_HEIGHT / 2) - 1,
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
		return []int{0, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 5}
	case 6:
		return []int{0, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 6}
	}

	return []int{0}
}

func (s *Station) GetObjects(direction int, fenceInside, fenceOutside bool) []properties.BoundingBox {
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

	baseSprites := 4
	if s.HasFences {
		baseSprites = 8
	}

	result := make([]properties.BoundingBox, 0)

	if fenceOutside {
		result = append(result, properties.BoundingBox{YOffset: yOffset, XOffset: xOffset, X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: s.GetBaseSpriteNumber() + base + 4})
	} else {
		result = append(result, properties.BoundingBox{YOffset: yOffset, XOffset: xOffset, X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: s.GetBaseSpriteNumber() + base + 0})
	}

	if fenceInside {
		result = append(result, properties.BoundingBox{X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: s.GetBaseSpriteNumber() + base + 5})
	} else {
		result = append(result, properties.BoundingBox{X: x, Y: y, Z: s.PlatformHeight, SpriteNumber: s.GetBaseSpriteNumber() + base + 1})
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
			SpriteNumber: s.GetBaseSpriteNumber() + baseSprites + (idx * 2) + direction,
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

	platformSprites := 4
	if s.HasFences {
		platformSprites = 8
	}
	file.AddElement(&Spritesets{ID: 0, NumSets: s.MaxLoadState + 1, NumSprites: platformSprites + (len(s.AdditionalObjects) * 2)})

	for i := 0; i <= s.MaxLoadState; i++ {
		filename := fmt.Sprintf("%s_%d_8bpp.png", s.SpriteFilename, i)

		// Non-fence sprites
		file.AddElement(&Sprites{
			GetSprite(filename, 0, false),
			GetSprite(filename, 1, false),
			GetSprite(filename, 2, true),
			GetSprite(filename, 3, true),
		})

		// Fence sprites
		if s.HasFences {
			fenceFilename := fmt.Sprintf("%s_%d_fence_8bpp.png", s.SpriteFilename, i)
			file.AddElement(&Sprites{
				GetSprite(fenceFilename, 0, false),
				GetSprite(fenceFilename, 1, false),
				GetSprite(fenceFilename, 2, true),
				GetSprite(fenceFilename, 3, true),
			})
		}

		for _, obj := range s.AdditionalObjects {
			filename := fmt.Sprintf("%s_%d_8bpp.png", obj.SpriteFilename, i)
			file.AddElement(&Sprites{
				GetSprite(filename, 0, obj.InvertDirection != true),
				GetSprite(filename, 1, obj.InvertDirection != false),
			})
		}
	}

	def := &Definition{StationID: s.ID}
	def.AddProperty(&properties.ClassID{ID: s.ClassID})

	if s.HasFences {
		def.AddProperty(&properties.CallbackFlag{SpriteLayout: true})
	}

	def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

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

		entry := properties.LayoutEntry{
			EastWest: properties.SpriteDirection{
				GroundSprite: 1012,
				Sprites:      s.GetObjects(EAST_WEST, fenceInside, fenceOutside),
			},
			NorthSouth: properties.SpriteDirection{
				GroundSprite: 1011,
				Sprites:      s.GetObjects(NORTH_SOUTH, fenceInside, fenceOutside),
			},
		}

		layoutEntries = append(layoutEntries, entry)
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

	if s.HasFences {
		file.AddElement(&callbacks.StationFenceCallback{
			SetID:            10,
			DefaultSpriteSet: 0,
		})

		file.AddElement(&callbacks.StationFenceCallback{
			SetID:            20,
			DefaultSpriteSet: 1,
		})

		passengerCargoSet, otherCargoSet = 10, 20
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
