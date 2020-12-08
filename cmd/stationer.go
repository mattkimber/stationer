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
	EmptySprite  int
	RoofSprite   int
	RoofPlatform int
}

const (
	// This is not the actual number, but the number leaving some room for expansion
	PLATFORM_TYPES  = 20
	CLASS_PLATFORMS = (PLATFORM_TYPES * 3) + 10
)

func main() {
	file := output_file.File{}
	file.AddElement(&nfo.Header{
		Initials:    "TWF",
		SetID:       8,
		SetName:     "Timberwolf's Stations 1.1.4",
		Description: "A set of British-style railway stations feature multiple eras of platforms, buildings and waypoints in 2x zoom",
		Version:     9,
		MinVersion:  4,
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
			Enable1: true,
		},
	}

	largeObjectConfiguration := properties.PlatformLayout{
		Platforms: properties.PlatformBitmask{
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
			Enable2:    true,
			Enable3:    true,
			Enable4:    true,
			Enable5:    true,
			Enable6:    true,
			Enable7:    true,
			EnableMore: true,
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
				{Filename: "bare_shelter_tiled_wall", HasFences: true, MaxLoadState: 5},
				{Filename: "ramp_ne", HasFences: true, MaxLoadState: 5},
				{Filename: "ramp_sw", HasFences: true, MaxLoadState: 5},
				{Filename: "bare_footbridge", HasFences: true, MaxLoadState: 5},
				{Filename: "billboard_1", HasFences: false, MaxLoadState: 5, DedicatedFlipSprite: true, SingleSided: true},
				{Filename: "billboard_2", HasFences: false, MaxLoadState: 5, DedicatedFlipSprite: true, SingleSided: true},
				{Filename: "billboard_3", HasFences: false, MaxLoadState: 5, DedicatedFlipSprite: true, SingleSided: true},
				{Filename: "billboard_4", HasFences: false, MaxLoadState: 5, DedicatedFlipSprite: true, SingleSided: true},
			},
			BaseFilename: class.Filename,
		}

		if class.ClassID != "TWF0" {
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "bare_footbridge_covered", HasFences: true, MaxLoadState: 5})
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "bare_footbridge_covered_brick", HasFences: true, MaxLoadState: 5})
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "bare_shelter_glass", HasFences: true, MaxLoadState: 5})
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "shelter_glass", HasFences: false, MaxLoadState: 5, IsStatic: true})
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "underpass", HasFences: true, MaxLoadState: 5})
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "bare_cafe", HasFences: false, MaxLoadState: 5})
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "bare_planter", HasFences: false, MaxLoadState: 5})
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "roof", HasFences: false, MaxLoadState: 5})
		}

		classSprites.SetStatistics()

		footbridgeSprites := []sprites.PlatformObject{
			{
				BaseSpriteID:   classSprites.LastSpriteNumber,
				SpriteFilename: "footbridge",
				MaxLoadState:   5,
			},
		}

		if class.ClassID != "TWF0" {
			footbridgeSprites = append(footbridgeSprites, []sprites.PlatformObject{
				{
					BaseSpriteID:   classSprites.LastSpriteNumber + 2,
					SpriteFilename: "footbridge_covered",
					MaxLoadState:   5,
					IsStatic:       true,
				},
				{
					BaseSpriteID:   classSprites.LastSpriteNumber + 4,
					SpriteFilename: "footbridge_covered_brick",
					MaxLoadState:   5,
					IsStatic:       true,
				},
			}...)
		}

		// +6 = footbridge sprites
		footbridgeSpriteCount := len(footbridgeSprites) * 2

		total := classSprites.LastSpriteNumber + footbridgeSpriteCount

		roofSprite := sprites.StationRoof{}
		if class.ClassID != "TWF0" {
			roofSprite = sprites.StationRoof{
				SpriteFilename: "",
				MaxLoadState:   5,
				RoofType:       "arch",
				BaseSpriteID:   classSprites.LastSpriteNumber + footbridgeSpriteCount,
			}

			class.EmptySprite = classSprites.SpriteMap["sign"]
			class.RoofPlatform = classSprites.SpriteMap["roof"]
			class.RoofSprite = roofSprite.BaseSpriteID

			// +12 = roof sprites
			total = total + 12
		}

		// Definition for all the spritesets
		file.AddElement(&sprites.Spritesets{ID: 0, NumSets: sprites.GLOBAL_MAX_LOAD_STATE + 1, NumSprites: total})

		// Write each type of sprite to the file
		for i := 0; i <= sprites.GLOBAL_MAX_LOAD_STATE; i++ {
			classSprites.WriteToFile(&file, i)
			for _, footbridgeSprite := range footbridgeSprites {
				footbridgeSprite.WriteToFile(&file, i)
			}

			if class.ClassID != "TWF0" {
				roofSprite.WriteToFile(&file, i)
			}
		}

		names := []string{"", "inner", "outer"}
		for i := 0; i < 3; i++ {
			baseObjectID := class.BaseObjectID + (PLATFORM_TYPES * i)
			inner := i <= 1
			outer := i == 0 || i == 2
			bracketName := ""
			commaName := ""

			if names[i] != "" {
				bracketName = " (" + names[i] + ")"
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
					YearAvailable:    max(class.Available, 1835),
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
					YearAvailable:    max(class.Available, 1845),
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
					ID:               baseObjectID + 14,
					BaseSpriteID:     classSprites.SpriteMap["bare_shelter_tiled_wall"],
					ClassID:          class.ClassID,
					ClassName:        class.ClassName,
					ObjectName:       "Shelter (with wall" + commaName + ")",
					YearAvailable:    max(class.Available, 1890),
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
					ObjectName:            "Ramp (NE" + commaName + ")",
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
					ObjectName:            "Ramp (SW" + commaName + ")",
					PlatformConfiguration: rampConfiguration,
					UseCompanyColour:      true,
					HasFences:             true,
					InnerPlatform:         inner,
					OuterPlatform:         outer,
				},
			}

			if class.ClassID != "TWF0" {
				shelterGlass := getGlassObjects(inner, outer, classSprites.SpriteMap["shelter_glass"], true)
				solarPanels := getGlassObjects(inner, outer, classSprites.SpriteMap["shelter_glass"], false)

				thisClass = append(thisClass, []nfo.Station{
					{
						ID:                    baseObjectID + 19,
						BaseSpriteID:          classSprites.SpriteMap["underpass"],
						ClassID:               class.ClassID,
						ClassName:             class.ClassName,
						PlatformConfiguration: rampConfiguration,
						ObjectName:            "Underpass" + bracketName,
						YearAvailable:         max(class.Available, 1890),
						MaxLoadState:          5,
						UseCompanyColour:      true,
						HasFences:             true,
						InnerPlatform:         inner,
						OuterPlatform:         outer,
					},
					{
						ID:                baseObjectID + 15,
						BaseSpriteID:      classSprites.SpriteMap["bare_shelter_glass"],
						ClassID:           class.ClassID,
						ClassName:         class.ClassName,
						ObjectName:        "Shelter (glass" + commaName + ")",
						YearAvailable:     max(class.Available, 1990),
						AdditionalObjects: shelterGlass,
						MaxLoadState:      5,
						UseCompanyColour:  true,
						HasFences:         true,
						InnerPlatform:     inner,
						OuterPlatform:     outer,
					},
					{
						ID:                baseObjectID + 16,
						BaseSpriteID:      classSprites.SpriteMap["bare_shelter_glass"],
						ClassID:           class.ClassID,
						ClassName:         class.ClassName,
						ObjectName:        "Shelter (solar panels" + commaName + ")",
						YearAvailable:     max(class.Available, 2005),
						AdditionalObjects: solarPanels,
						MaxLoadState:      5,
						UseCompanyColour:  true,
						HasFences:         true,
						InnerPlatform:     inner,
						OuterPlatform:     outer,
					},
				}...)
			}

			// Only available if we have an "inner" platform
			if i <= 1 {

				thisClass = append(thisClass, []nfo.Station{
					{
						ID:           baseObjectID + 7,
						BaseSpriteID: classSprites.SpriteMap["billboard_1"],
						RandomSpriteIDs: []int{
							classSprites.SpriteMap["billboard_1"],
							classSprites.SpriteMap["billboard_2"],
							classSprites.SpriteMap["billboard_3"],
							classSprites.SpriteMap["billboard_4"],
						},
						ClassID:               class.ClassID,
						ClassName:             class.ClassName,
						YearAvailable:         max(class.Available, 1845),
						MaxLoadState:          5,
						ObjectName:            "Billboard" + bracketName,
						UseCompanyColour:      true,
						HasFences:             true,
						InnerPlatform:         inner,
						OuterPlatform:         outer,
						HasLargeCentralObject: true,
						OverrideOuter:         true,
						ObjectIsSingleSided:   true,
						OuterPlatformSprite:   classSprites.SpriteMap["empty"],
						PlatformConfiguration: largeObjectConfiguration,
					},
				}...)

				// Tiles only available on concrete/modern
				if class.ClassID != "TWF0" {
					thisClass = append(thisClass, []nfo.Station{
						{
							ID:                    baseObjectID + 8,
							BaseSpriteID:          classSprites.SpriteMap["bare_cafe"],
							ClassID:               class.ClassID,
							ClassName:             class.ClassName,
							YearAvailable:         max(class.Available, 1932),
							MaxLoadState:          5,
							ObjectName:            "Waiting Room" + bracketName,
							UseCompanyColour:      true,
							HasFences:             true,
							InnerPlatform:         inner,
							OuterPlatform:         outer,
							HasLargeCentralObject: true,
							OverrideOuter:         true,
							OuterPlatformSprite:   classSprites.SpriteMap["empty"],
							PlatformConfiguration: largeObjectConfiguration,
						},
						{
							ID:                    baseObjectID + 9,
							BaseSpriteID:          classSprites.SpriteMap["bare_planter"],
							ClassID:               class.ClassID,
							ClassName:             class.ClassName,
							YearAvailable:         max(class.Available, 1870),
							MaxLoadState:          5,
							ObjectName:            "Planter" + bracketName,
							UseCompanyColour:      true,
							HasFences:             true,
							InnerPlatform:         inner,
							OuterPlatform:         outer,
							HasLargeCentralObject: true,
							OverrideOuter:         true,
							OuterPlatformSprite:   classSprites.SpriteMap["empty"],
							PlatformConfiguration: largeObjectConfiguration,
						},
					}...)
				}
			}

			if i == 0 {
				thisClass = append(thisClass, []nfo.Station{
					{
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
								BaseSpriteID: footbridgeSprites[0].BaseSpriteID,
							},
						},
					}}...)

				if class.ClassID != "TWF0" {
					thisClass = append(thisClass, []nfo.Station{
						{
							ID:                    baseObjectID + 17,
							BaseSpriteID:          classSprites.SpriteMap["bare_footbridge_covered_brick"],
							ClassID:               class.ClassID,
							ClassName:             class.ClassName,
							YearAvailable:         max(class.Available, 1910),
							ObjectName:            "Footbridge (covered)",
							UseCompanyColour:      true,
							HasFences:             true,
							MaxLoadState:          5,
							PlatformHeight:        16,
							InnerPlatform:         true,
							OuterPlatform:         true,
							PlatformConfiguration: rampConfiguration,
							AdditionalObjects: []nfo.AdditionalObject{
								{
									X:            6,
									Y:            2,
									Z:            15,
									SizeX:        5,
									SizeY:        8,
									SizeZ:        3,
									BaseSpriteID: footbridgeSprites[2].BaseSpriteID,
								},
							},
						},
						{
							ID:                    baseObjectID + 18,
							BaseSpriteID:          classSprites.SpriteMap["bare_footbridge_covered"],
							ClassID:               class.ClassID,
							ClassName:             class.ClassName,
							YearAvailable:         max(class.Available, 1932),
							ObjectName:            "Footbridge (covered)",
							UseCompanyColour:      true,
							HasFences:             true,
							MaxLoadState:          5,
							PlatformHeight:        16,
							InnerPlatform:         true,
							OuterPlatform:         true,
							PlatformConfiguration: rampConfiguration,
							AdditionalObjects: []nfo.AdditionalObject{
								{
									X:            6,
									Y:            2,
									Z:            15,
									SizeX:        5,
									SizeY:        8,
									SizeZ:        3,
									BaseSpriteID: footbridgeSprites[1].BaseSpriteID,
								},
							},
						},
					}...)
				}

				thisClass = append(thisClass, []nfo.Station{
					{
						ID:                    baseObjectID + 10,
						BaseSpriteID:          classSprites.SpriteMap["ramp_ne"],
						ClassID:               class.ClassID,
						ClassName:             class.ClassName,
						YearAvailable:         class.Available,
						MaxLoadState:          5,
						ObjectName:            "Half Ramp (NE, inner" + commaName + ")",
						UseCompanyColour:      true,
						HasFences:             true,
						InnerPlatform:         inner,
						OuterPlatform:         outer,
						OverrideOuter:         true,
						OuterPlatformSprite:   classSprites.SpriteMap["empty"],
						PlatformConfiguration: rampConfiguration,
					},
					{
						ID:                    baseObjectID + 11,
						BaseSpriteID:          classSprites.SpriteMap["empty"],
						ClassID:               class.ClassID,
						ClassName:             class.ClassName,
						YearAvailable:         class.Available,
						MaxLoadState:          5,
						ObjectName:            "Half Ramp (NE, outer" + commaName + ")",
						UseCompanyColour:      true,
						HasFences:             true,
						InnerPlatform:         inner,
						OuterPlatform:         outer,
						OverrideOuter:         true,
						OuterPlatformSprite:   classSprites.SpriteMap["ramp_ne"],
						PlatformConfiguration: rampConfiguration,
					},
					{
						ID:                    baseObjectID + 12,
						BaseSpriteID:          classSprites.SpriteMap["ramp_sw"],
						ClassID:               class.ClassID,
						ClassName:             class.ClassName,
						YearAvailable:         class.Available,
						MaxLoadState:          5,
						ObjectName:            "Half Ramp (SW, inner" + commaName + ")",
						UseCompanyColour:      true,
						HasFences:             true,
						InnerPlatform:         inner,
						OuterPlatform:         outer,
						OverrideOuter:         true,
						OuterPlatformSprite:   classSprites.SpriteMap["empty"],
						PlatformConfiguration: rampConfiguration,
					},
					{
						ID:                    baseObjectID + 13,
						BaseSpriteID:          classSprites.SpriteMap["empty"],
						ClassID:               class.ClassID,
						ClassName:             class.ClassName,
						YearAvailable:         class.Available,
						MaxLoadState:          5,
						ObjectName:            "Half Ramp (SW, outer" + commaName + ")",
						UseCompanyColour:      true,
						HasFences:             true,
						InnerPlatform:         inner,
						OuterPlatform:         outer,
						OverrideOuter:         true,
						OuterPlatformSprite:   classSprites.SpriteMap["ramp_sw"],
						PlatformConfiguration: rampConfiguration,
					},
				}...)
			}

			for _, station := range thisClass {
				station.WriteToFile(&file)
			}
		}

		// Wooden platforms do not have station halls
		if class.ClassID != "TWF0" {

			hall := nfo.StationHall{
				ID:                    class.BaseObjectID + 31,
				BarePlatformSprite:    class.EmptySprite,
				RoofPlatformSprite:    class.RoofPlatform,
				RoofBaseSprite:        class.RoofSprite,
				ClassID:               class.ClassID,
				ClassName:             class.ClassName,
				ObjectName:            "Station Hall",
				PlatformConfiguration: properties.PlatformLayout{},
				UseCompanyColour:      true,
				MaxLoadState:          5,
				PlatformHeight:        0,
				YearAvailable:         max(class.Available, 1870),
			}

			hall.WriteToFile(&file)
		}

		buffers := nfo.BufferStop{
			ID:               class.BaseObjectID + 32,
			SpriteFilename:   fmt.Sprintf("%s_bufferstop", class.Filename),
			ClassID:          class.ClassID,
			ClassName:        class.ClassName,
			YearAvailable:    class.Available,
			ObjectName:       "Buffer Stop",
			UseCompanyColour: true,
		}

		buffers.WriteToFile(&file)

		platforms := []nfo.FullTilePlatform{
			{
				SpriteFilename:   fmt.Sprintf("%s_concourse", class.Filename),
				ID:               class.BaseObjectID + (PLATFORM_TYPES * 3) + 3,
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				YearAvailable:    class.Available,
				ObjectName:       "Concourse",
				UseCompanyColour: true,
			},
			{
				SpriteFilename:   fmt.Sprintf("%s_concourse_shelter", class.Filename),
				ID:               class.BaseObjectID + (PLATFORM_TYPES * 3) + 4,
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				YearAvailable:    max(class.Available, 1860),
				ObjectName:       "Concourse with shelters",
				UseCompanyColour: true,
			},
		}

		for _, platform := range platforms {
			platform.WriteToFile(&file)
		}

	}

	objectID := CLASS_PLATFORMS * 3

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

func getGlassObjects(inner, outer bool, baseSpriteID int, isTransparent bool) []nfo.AdditionalObject {
	shelterGlass := make([]nfo.AdditionalObject, 0)

	if inner {
		shelterGlass = append(shelterGlass, nfo.AdditionalObject{
			X:                0,
			Y:                0,
			Z:                0,
			SizeX:            16,
			SizeY:            5,
			SizeZ:            5,
			IsTransparent:    isTransparent,
			BaseSpriteID:     baseSpriteID,
			HasFourWaySprite: true,
		})
	}

	if outer {
		shelterGlass = append(shelterGlass, nfo.AdditionalObject{
			X:                0,
			Y:                16 - 5,
			Z:                0,
			SizeX:            16,
			SizeY:            5,
			SizeZ:            5,
			IsTransparent:    isTransparent,
			BaseSpriteID:     baseSpriteID + 1,
			HasFourWaySprite: true,
		})
	}
	return shelterGlass
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
