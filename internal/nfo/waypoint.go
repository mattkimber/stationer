package nfo

import (
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo/callbacks"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
	"github.com/mattkimber/stationer/internal/nfo/sprites"
)

type Waypoint struct {
	ID               int
	SpriteFilename   string
	ClassID          string
	ClassName        string
	ObjectName       string
	UseCompanyColour bool
	YearAvailable    int
}

const (
	WAYPOINT_SPRITE_HEIGHT      = 55
	WAYPOINT_BASE_SPRITE_HEIGHT = 35
)

func (wp *Waypoint) SetID(id int) {
	wp.ID = id
}

func (wp *Waypoint) GetID() int {
	return wp.ID
}

func GetWaypointSprite(filename string, num int, swap bool) sprites.Sprite {
	xrel := -(BUILDING_SPRITE_WIDTH / 2) - 10
	yrel := -(WAYPOINT_SPRITE_HEIGHT / 2) - 1

	if swap {
		xrel = 11 - (BUILDING_SPRITE_WIDTH / 2)
	}

	return sprites.Sprite{
		Filename: filename,
		X:        BUILDING_SPRITE_WIDTH_WITH_PADDING * num,
		Y:        0,
		Width:    BUILDING_SPRITE_WIDTH,
		Height:   WAYPOINT_SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     yrel,
	}
}

func GetWaypointBaseSprite(filename string, num int, swap bool) sprites.Sprite {
	xrel := 1 - (BUILDING_SPRITE_WIDTH / 2)
	yrel := -(WAYPOINT_BASE_SPRITE_HEIGHT / 2) + 14

	if swap {
		xrel = -(BUILDING_SPRITE_WIDTH / 2)
	}

	return sprites.Sprite{
		Filename: filename,
		X:        BUILDING_SPRITE_WIDTH_WITH_PADDING * num,
		Y:        0,
		Width:    BUILDING_SPRITE_WIDTH,
		Height:   WAYPOINT_BASE_SPRITE_HEIGHT,
		XRel:     xrel,
		YRel:     yrel,
	}
}

func (wp *Waypoint) GetBaseSpriteNumber() int {
	if wp.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE
	}

	return CUSTOM_SPRITE
}

func (wp *Waypoint) GetObjects(direction int, withBuilding bool) []properties.BoundingBox {
	x, y := 16, 5

	if direction == NORTH_SOUTH {
		x = 5
		y = 16
	}

	result := make([]properties.BoundingBox, 0)
	base := 0
	if direction == NORTH_SOUTH {
		base = 1
	}

	if withBuilding {
		// Append the building and half base sprite
		result = append(result, properties.BoundingBox{X: 16, Y: 16, Z: 1, SpriteNumber: wp.GetBaseSpriteNumber() + base + 2})
		result = append(result, properties.BoundingBox{X: x, Y: y, Z: 12, SpriteNumber: wp.GetBaseSpriteNumber() + base})
	} else {
		// Append the full base sprite
		result = append(result, properties.BoundingBox{X: 16, Y: 16, Z: 1, SpriteNumber: wp.GetBaseSpriteNumber() + base + 4})
	}

	return result
}

func (wp *Waypoint) WriteToFile(file *output_file.File) {
	wp.addSprites(file)

	def := &Definition{StationID: wp.ID}
	def.AddProperty(&properties.ClassID{ID: wp.ClassID})

	// This is irrelevant?
	// def.AddProperty(&properties.LittleLotsThreshold{Amount: 200})

	layoutEntries := make([]properties.LayoutEntry, 0)

	// Add the layouts (2 layouts, 1 with building, 1 without)
	layoutEntries = append(layoutEntries, wp.getLayoutEntry(true))
	layoutEntries = append(layoutEntries, wp.getLayoutEntry(false))

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

func (wp *Waypoint) addSprites(file *output_file.File) {
	file.AddElement(&sprites.Spritesets{ID: 0, NumSets: 1, NumSprites: 6})

	// Waypoint building
	filename := fmt.Sprintf("%s_8bpp.png", wp.SpriteFilename)

	file.AddElement(&sprites.Sprites{
		GetWaypointSprite(filename, 0, false),
		GetWaypointSprite(filename, 1, true),
	})

	// Waypoint crossing (for building tile)
	filename = fmt.Sprintf("%s_base_8bpp.png", wp.SpriteFilename)

	file.AddElement(&sprites.Sprites{
		GetWaypointBaseSprite(filename, 0, false),
		GetWaypointBaseSprite(filename, 1, true),
	})

	// Waypoint crossing - full
	filename = fmt.Sprintf("%s_base_full_8bpp.png", wp.SpriteFilename)

	file.AddElement(&sprites.Sprites{
		GetWaypointBaseSprite(filename, 0, false),
		GetWaypointBaseSprite(filename, 1, true),
	})

}

func (wp *Waypoint) getLayoutEntry(withBuilding bool) properties.LayoutEntry {
	entry := properties.LayoutEntry{
		EastWest: properties.SpriteDirection{
			GroundSprite: GROUND_SPRITE_RAIL_EW,
			// East-West sprites are assembled in the "wrong" order so that
			// multi tile stations are the correct way round when displayed
			Sprites: wp.GetObjects(EAST_WEST, withBuilding),
		},
		NorthSouth: properties.SpriteDirection{
			GroundSprite: GROUND_SPRITE_RAIL_NS,
			Sprites:      wp.GetObjects(NORTH_SOUTH, withBuilding),
		},
	}
	return entry
}
