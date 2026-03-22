package nfo

import (
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
	"github.com/mattkimber/stationer/internal/nfo/sprites"
)

type BridgeWaypoint struct {
	ID               int
	ClassID          string
	ClassName        string
	ObjectName       string
	UseCompanyColour bool
	YearAvailable    int
	AdditionalSprites     []sprites.PlatformObject
	AdditionalObjects     []AdditionalObject
}

func (wp *BridgeWaypoint) SetID(id int) {
	wp.ID = id
}

func (wp *BridgeWaypoint) GetID() int {
	return wp.ID
}

func (wp *BridgeWaypoint) GetBaseSpriteNumber() int {
	if wp.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func (wp *BridgeWaypoint) GetObjects(direction int) []properties.BoundingBox {
	result := make([]properties.BoundingBox, 0)
	
	for _, obj := range wp.AdditionalObjects {
		x, y := obj.SizeX, obj.SizeY
		xOffset, yOffset := obj.X, obj.Y

		if direction == NORTH_SOUTH {
			x, y = obj.SizeY, obj.SizeX
			xOffset, yOffset = obj.Y, obj.X
		}

		multiplier := 1
		if obj.HasFourWaySprite {
			multiplier = 2
		}

		result = append(result, properties.BoundingBox{
			XOffset:      xOffset,
			YOffset:      yOffset,
			ZOffset:      obj.Z,
			X:            x,
			Y:            y,
			Z:            obj.SizeZ,
			SpriteNumber: wp.GetBaseSpriteNumber() + obj.BaseSpriteID + (direction * multiplier),
		})
	}

	return result
}

func (wp *BridgeWaypoint) WriteToFile(file *output_file.File) {
	wp.addSprites(file)

	def := &Definition{StationID: wp.ID}
	def.AddProperty(&properties.ClassID{ID: wp.ClassID})

    // There's only two bridge clearance entries needed, since we have one layout entry with two orientations.
	def.AddProperty(&properties.MinimumBridgeClearance{Clearance: 2, Layouts: 2})

	// This is irrelevant?
	// def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	layoutEntries := make([]properties.LayoutEntry, 0)

	// There's only one layout entry for a bridge.
	layoutEntries = append(layoutEntries, wp.getLayoutEntry())

	def.AddProperty(&properties.SpriteLayout{
		Entries: layoutEntries,
	})

	def.AddProperty(&properties.CallbackFlag{SpriteLayout: true, Availability: wp.YearAvailable != 0})

	file.AddElement(def)

	file.AddElement(&StationSet{
		SetID:         0,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	spriteset, yearCallbackID := 2, 0

	if wp.YearAvailable != 0 {
		yearCallbackID = 2
		file.AddElement(&callbacks.AvailabilityYearCallback{
			SetID:      yearCallbackID,
			HasDecider: false,
			Year:       wp.YearAvailable,
		})
	}

	// Add the callback for building tile selection
	file.AddElement(&callbacks.WaypointSpriteCallback{YearCallbackID: yearCallbackID})

	file.AddElement(&GraphicSetAssignment{
		IDs:        []int{wp.ID},
		DefaultSet: spriteset,
	})

	file.AddElement(&TextString{
		LanguageFile:   255,
		StationId:      wp.ID,
		TextStringType: TextStringTypeStationName,
		Text:           wp.ObjectName,
	})

	file.AddElement(&TextString{
		LanguageFile:   255,
		StationId:      wp.ID,
		TextStringType: TextStringTypeClassName,
		Text:           wp.ClassName,
	})
}

func (wp *BridgeWaypoint) addSprites(file *output_file.File) {
	file.AddElement(&sprites.Spritesets{ID: 0, NumSets: 1, NumSprites: len(wp.AdditionalObjects) * 2})

    for _, obj := range wp.AdditionalSprites {
        obj.WriteToFile(file, 0)
    }

}

func (wp *BridgeWaypoint) getLayoutEntry() properties.LayoutEntry {
	entry := properties.LayoutEntry{
		EastWest: properties.SpriteDirection{
			GroundSprite: GROUND_SPRITE_RAIL_EW,
			// East-West sprites are assembled in the "wrong" order so that
			// multi tile stations are the correct way round when displayed
			Sprites: wp.GetObjects(EAST_WEST),
		},
		NorthSouth: properties.SpriteDirection{
			GroundSprite: GROUND_SPRITE_RAIL_NS,
			Sprites:      wp.GetObjects(NORTH_SOUTH),
		},
	}
	return entry
}
