package main

import (
	"flag"
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
	"github.com/mattkimber/stationer/internal/nfo/sprites"
)

func init() {
	flag.Parse()
}

type StationClass struct {
	Filename     string
	ClassID      string
	ClassName    string
	Available    int
	BaseObjectID int
}

const (
	// This is not the actual number, but the number leaving some room for expansion
	PLATFORM_TYPES  = 10
	CLASS_PLATFORMS = PLATFORM_TYPES * 3
)

func main() {
	file := output_file.File{}
	file.AddElement(&nfo.Header{
		Initials:    "TWF",
		SetID:       8,
		SetName:     "Timberwolf's Stations 1.0.2",
		Description: "A set of British-style railway stations feature multiple eras of platforms, buildings and waypoints in 2x zoom",
		Version:     3,
		MinVersion:  2,
	})

	file.AddElement(&nfo.CargoTypeTable{Cargos: []string{"PASS", "MAIL"}})

	classes := []StationClass{
		{Filename: "wooden", ClassID: "TWF0", ClassName: "Wooden Platforms", Available: 0, BaseObjectID: CLASS_PLATFORMS * 0},
		{Filename: "concrete", ClassID: "TWF1", ClassName: "Concrete Platforms", Available: 1860, BaseObjectID: CLASS_PLATFORMS * 1},
		{Filename: "modern", ClassID: "TWF2", ClassName: "Modern Platforms", Available: 1970, BaseObjectID: CLASS_PLATFORMS * 2},
	}

	rampConfiguration := properties.PlatformLayout{
		Platforms: properties.PlatformBitmask{
			Enable1:    true,
			Enable2:    true,
			Enable3:    true,
			Enable4:    true,
			Enable5:    true,
			Enable6:    true,
			Enable7:    true,
			EnableMore: true,
		},
		Lengths: properties.PlatformBitmask{
			Enable1:    true,
			Enable2:    false,
			Enable3:    false,
			Enable4:    false,
			Enable5:    false,
			Enable6:    false,
			Enable7:    false,
			EnableMore: false,
		},
	}

	for _, class := range classes {

		// All the station sprites are put in one massive spriteset
		// so they can be mixed and matched
		classSprites := sprites.StationSprites{
			Sprites: []sprites.StationSprite{
				{Filename: "empty", HasFences: true, MaxLoadState: 6},
				{Filename: "sign", HasFences: true, MaxLoadState: 6},
				{Filename: "benches", HasFences: true, MaxLoadState: 6},
				{Filename: "bare_shelter_traditional", HasFences: true, MaxLoadState: 5},
				{Filename: "ramp_ne", HasFences: true, MaxLoadState: 5},
				{Filename: "ramp_sw", HasFences: true, MaxLoadState: 5},
				{Filename: "bare_footbridge", HasFences: true, MaxLoadState: 5},
			},
			BaseFilename: class.Filename,
		}

		classSprites.SetStatistics()

		footbridgeSprite := sprites.PlatformObject{
			BaseSpriteID:   classSprites.LastSpriteNumber,
			SpriteFilename: "footbridge",
			MaxLoadState:   5,
		}

		// +2 = footbridge sprite
		total := classSprites.LastSpriteNumber + 2

		// Definition for all the spritesets
		file.AddElement(&sprites.Spritesets{ID: 0, NumSets: sprites.GLOBAL_MAX_LOAD_STATE + 1, NumSprites: total})

		// Write each type of sprite to the file
		for i := 0; i <= sprites.GLOBAL_MAX_LOAD_STATE; i++ {
			classSprites.WriteToFile(&file, i)
			footbridgeSprite.WriteToFile(&file, i)
		}

		names := []string{"", "inner", "outer"}
		for i := 0; i < 3; i++ {
			baseObjectID := class.BaseObjectID + (PLATFORM_TYPES * i)
			inner := i <= 1
			outer := i == 0 || i == 2
			bracketName := ""
			commaName := ""

			if names[i] != "" {
				bracketName = "(" + names[i] + ")"
				commaName = ", " + names[i]
			}

			thisClass := []nfo.Station{
				{
					ID:               baseObjectID + 0,
					BaseSpriteID:     classSprites.SpriteMap["empty"],
					ClassID:          class.ClassID,
					ClassName:        class.ClassName,
					ObjectName:       "Platform" + bracketName,
					YearAvailable:    class.Available,
					UseCompanyColour: true,
					HasFences:        true,
					InnerPlatform:    inner,
					OuterPlatform:    outer,
				},
				{
					ID:               baseObjectID + 1,
					BaseSpriteID:     classSprites.SpriteMap["sign"],
					ClassID:          class.ClassID,
					ClassName:        class.ClassName,
					ObjectName:       "Platform with sign" + bracketName,
					YearAvailable:    max(class.Available, 1845),
					UseCompanyColour: true,
					HasFences:        true,
					InnerPlatform:    inner,
					OuterPlatform:    outer,
				},
				{
					ID:               baseObjectID + 2,
					BaseSpriteID:     classSprites.SpriteMap["benches"],
					ClassID:          class.ClassID,
					ClassName:        class.ClassName,
					ObjectName:       "Platform with benches" + bracketName,
					YearAvailable:    max(class.Available, 1840),
					UseCompanyColour: true,
					HasFences:        true,
					InnerPlatform:    inner,
					OuterPlatform:    outer,
				},
				{
					ID:               baseObjectID + 3,
					BaseSpriteID:     classSprites.SpriteMap["bare_shelter_traditional"],
					ClassID:          class.ClassID,
					ClassName:        class.ClassName,
					ObjectName:       "Shelter (traditional" + commaName + ")",
					YearAvailable:    max(class.Available, 1860),
					MaxLoadState:     5,
					UseCompanyColour: true,
					HasFences:        true,
					InnerPlatform:    inner,
					OuterPlatform:    outer,
				},
				{
					ID:                    baseObjectID + 4,
					BaseSpriteID:          classSprites.SpriteMap["ramp_ne"],
					ClassID:               class.ClassID,
					ClassName:             class.ClassName,
					YearAvailable:         class.Available,
					MaxLoadState:          5,
					ObjectName:            "Ramp (NE)" + commaName + ")",
					PlatformConfiguration: rampConfiguration,
					UseCompanyColour:      true,
					HasFences:             true,
					InnerPlatform:         inner,
					OuterPlatform:         outer,
				},
				{
					ID:                    baseObjectID + 5,
					BaseSpriteID:          classSprites.SpriteMap["ramp_sw"],
					ClassID:               class.ClassID,
					ClassName:             class.ClassName,
					YearAvailable:         class.Available,
					MaxLoadState:          5,
					ObjectName:            "Ramp (SW)" + commaName + ")",
					PlatformConfiguration: rampConfiguration,
					UseCompanyColour:      true,
					HasFences:             true,
					InnerPlatform:         inner,
					OuterPlatform:         outer,
				},
			}

			if i == 0 {
				thisClass = append(thisClass, nfo.Station{
					ID:                    baseObjectID + 6,
					BaseSpriteID:          classSprites.SpriteMap["bare_footbridge"],
					ClassID:               class.ClassID,
					ClassName:             class.ClassName,
					YearAvailable:         max(class.Available, 1865),
					ObjectName:            "Footbridge",
					UseCompanyColour:      true,
					HasFences:             true,
					MaxLoadState:          5,
					PlatformHeight:        16,
					InnerPlatform:         true,
					OuterPlatform:         true,
					PlatformConfiguration: rampConfiguration,
					AdditionalObjects: []nfo.AdditionalObject{
						{
							X:            5,
							Y:            4,
							Z:            13,
							SizeX:        5,
							SizeY:        8,
							SizeZ:        3,
							BaseSpriteID: footbridgeSprite.BaseSpriteID,
						},
					},
				})
			}

			for _, station := range thisClass {
				station.WriteToFile(&file)
			}
		}

	}

	objectID := 90

	// TODO: clean up and integrate bufferstops and station roofs properly
	for _, class := range classes {

		// Wooden platforms do not have station halls
		if class.ClassID != "TWF0" {

			hall := nfo.StationHall{
				ID:               objectID,
				SpriteFilename:   fmt.Sprintf("%s_empty", class.Filename),
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				YearAvailable:    max(class.Available, 1870),
				MaxLoadState:     5,
				ObjectName:       "Station Hall",
				RoofType:         "arch",
				UseCompanyColour: true,
			}

			hall.WriteToFile(&file)
			objectID = objectID + 1

		}

		buffers := nfo.BufferStop{
			ID:               objectID,
			SpriteFilename:   fmt.Sprintf("%s_bufferstop", class.Filename),
			ClassID:          class.ClassID,
			ClassName:        class.ClassName,
			YearAvailable:    class.Available,
			ObjectName:       "Buffer Stop",
			UseCompanyColour: true,
		}

		buffers.WriteToFile(&file)
		objectID = objectID + 1

		platforms := []nfo.FullTilePlatform{
			{
				SpriteFilename:   fmt.Sprintf("%s_concourse", class.Filename),
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				YearAvailable:    class.Available,
				ObjectName:       "Concourse",
				UseCompanyColour: true,
			},
			{
				SpriteFilename:   fmt.Sprintf("%s_concourse_shelter", class.Filename),
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				YearAvailable:    max(class.Available, 1860),
				ObjectName:       "Concourse with shelters",
				UseCompanyColour: true,
			},
		}

		for _, platform := range platforms {
			platform.ID = objectID
			platform.WriteToFile(&file)
			objectID = objectID + 1
		}
	}

	buildings := []nfo.Building{
		{
			SpriteFilename:   "wooden",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Wooden Station",
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "rural",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Rural Station",
			YearAvailable:    1840,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "art_deco",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Art Deco Station",
			YearAvailable:    1935,
			Width:            2,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "suburban_flat_roof",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Suburban Flat Roof Station",
			YearAvailable:    1962,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "small_city",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "City Station",
			YearAvailable:    1965,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "leslie_green",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Leslie Green Station",
			YearAvailable:    1904,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "j_m_easton",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "J.M. Easton Station",
			YearAvailable:    1935,
			UseCompanyColour: true,
		},
	}

	for _, building := range buildings {
		building.ID = objectID
		building.WriteToFile(&file)
		objectID = objectID + 1
	}

	waypoints := []nfo.Waypoint{
		{
			SpriteFilename:   "wp_kilby_bridge",
			ClassID:          "WAYP",
			ClassName:        "Waypoints",
			ObjectName:       "Kilby Bridge",
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "wp_hebden_bridge",
			ClassID:          "WAYP",
			ClassName:        "Waypoints",
			ObjectName:       "Hebden Bridge",
			YearAvailable:    1890,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "wp_sr_type_13",
			ClassID:          "WAYP",
			ClassName:        "Waypoints",
			ObjectName:       "SR Type 13",
			YearAvailable:    1930,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "wp_arnside",
			ClassID:          "WAYP",
			ClassName:        "Waypoints",
			ObjectName:       "Arnside",
			YearAvailable:    1870,
			UseCompanyColour: true,
		},
	}

	for _, waypoint := range waypoints {
		waypoint.ID = objectID
		waypoint.WriteToFile(&file)
		objectID = objectID + 1
	}

	file.Output()
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
