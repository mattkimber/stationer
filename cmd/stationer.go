package main

import (
	"flag"
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo"
	"github.com/mattkimber/stationer/internal/nfo/output_file"
	"github.com/mattkimber/stationer/internal/nfo/properties"
	"github.com/mattkimber/stationer/internal/nfo/sprites"
	"sort"
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
	PLATFORM_TYPES  = 30
	CLASS_PLATFORMS = (PLATFORM_TYPES * 3) + 10
)

func main() {
	outputObjects := make([]output_file.SortableFileWriter, 0)
	currentID := 0

	file := output_file.File{}
	file.AddElement(&nfo.Header{
		Initials:    "TWF",
		SetID:       8,
		SetName:     "Timberwolf's Stations 1.2.2",
		Description: "A set of British-style railway stations feature multiple eras of platforms, buildings and waypoints in 2x zoom",
		Version:     14,
		MinVersion:  12,
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
				{Filename: "fence", MaxLoadState: 5, IsStatic: true, IsRailFence: true, SingleSided: true},
				{Filename: "empty", HasFences: true, MaxLoadState: 6},
				{Filename: "sign", HasFences: true, MaxLoadState: 6},
				{Filename: "benches", HasFences: true, MaxLoadState: 6},
				{Filename: "bare_shelter_traditional", HasFences: true, MaxLoadState: 5},
				{Filename: "bare_shelter_tiled_wall", HasFences: true, MaxLoadState: 5},
				{Filename: "bare_hut", HasFences: true, MaxLoadState: 5},
				{Filename: "bare_stairs", HasFences: true, MaxLoadState: 5},
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
			classSprites.Sprites = append(classSprites.Sprites, sprites.StationSprite{Filename: "bare_shelter_curved", HasFences: true, MaxLoadState: 5})
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
				{
					BaseSpriteID:   classSprites.LastSpriteNumber + 6,
					SpriteFilename: "footbridge_pillar_covered",
					MaxLoadState:   5,
					IsStatic:       true,
				},
				{
					BaseSpriteID:   classSprites.LastSpriteNumber + 8,
					SpriteFilename: "footbridge_pillar_covered_b",
					MaxLoadState:   5,
					IsStatic:       true,
				},
				{
					BaseSpriteID:   classSprites.LastSpriteNumber + 10,
					SpriteFilename: "footbridge_pillar_covered_brick",
					MaxLoadState:   5,
					IsStatic:       true,
				},
				{
					BaseSpriteID:   classSprites.LastSpriteNumber + 12,
					SpriteFilename: "footbridge_pillar_covered_brick_b",
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
					ID:                baseObjectID + 0,
					BaseSpriteID:      classSprites.SpriteMap["empty"],
					ClassID:           class.ClassID,
					ClassName:         class.ClassName,
					ObjectName:        "Platform" + bracketName,
					YearAvailable:     class.Available,
					UseCompanyColour:  true,
					HasFences:         true,
					InnerPlatform:     inner,
					OuterPlatform:     outer,
					RailFenceSpriteID: classSprites.SpriteMap["fence"],
				},
				{
					ID:                baseObjectID + 1,
					BaseSpriteID:      classSprites.SpriteMap["sign"],
					ClassID:           class.ClassID,
					ClassName:         class.ClassName,
					ObjectName:        "Platform with sign" + bracketName,
					YearAvailable:     max(class.Available, 1835),
					UseCompanyColour:  true,
					HasFences:         true,
					InnerPlatform:     inner,
					OuterPlatform:     outer,
					RailFenceSpriteID: classSprites.SpriteMap["fence"],
				},
				{
					ID:                baseObjectID + 2,
					BaseSpriteID:      classSprites.SpriteMap["benches"],
					ClassID:           class.ClassID,
					ClassName:         class.ClassName,
					ObjectName:        "Platform with benches" + bracketName,
					YearAvailable:     max(class.Available, 1845),
					UseCompanyColour:  true,
					HasFences:         true,
					InnerPlatform:     inner,
					OuterPlatform:     outer,
					RailFenceSpriteID: classSprites.SpriteMap["fence"],
				},
				{
					ID:                baseObjectID + 3,
					BaseSpriteID:      classSprites.SpriteMap["bare_shelter_traditional"],
					ClassID:           class.ClassID,
					ClassName:         class.ClassName,
					ObjectName:        "Shelter (traditional" + commaName + ")",
					YearAvailable:     max(class.Available, 1860),
					MaxLoadState:      5,
					UseCompanyColour:  true,
					HasFences:         true,
					InnerPlatform:     inner,
					OuterPlatform:     outer,
					RailFenceSpriteID: classSprites.SpriteMap["fence"],
				},
				{
					ID:                baseObjectID + 14,
					BaseSpriteID:      classSprites.SpriteMap["bare_shelter_tiled_wall"],
					ClassID:           class.ClassID,
					ClassName:         class.ClassName,
					ObjectName:        "Shelter (with wall" + commaName + ")",
					YearAvailable:     max(class.Available, 1890),
					MaxLoadState:      5,
					UseCompanyColour:  true,
					HasFences:         true,
					InnerPlatform:     inner,
					OuterPlatform:     outer,
					RailFenceSpriteID: classSprites.SpriteMap["fence"],
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
					RailFenceSpriteID:     classSprites.SpriteMap["fence"],
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
					RailFenceSpriteID:     classSprites.SpriteMap["fence"],
				},
				{
					ID:                    baseObjectID + 22,
					BaseSpriteID:          classSprites.SpriteMap["bare_stairs"],
					ClassID:               class.ClassID,
					ClassName:             class.ClassName,
					YearAvailable:         max(class.Available, 1860),
					MaxLoadState:          5,
					ObjectName:            "Stairs" + bracketName,
					PlatformConfiguration: rampConfiguration,
					UseCompanyColour:      true,
					HasFences:             true,
					InnerPlatform:         inner,
					OuterPlatform:         outer,
					RailFenceSpriteID:     classSprites.SpriteMap["fence"],
				},
			}

			if inner && outer {
				thisClass = append(thisClass, []nfo.Station{{
				ID:                    baseObjectID + 23,
					BaseSpriteID:          classSprites.SpriteMap["bare_stairs"],
					ClassID:               class.ClassID,
					ClassName:             class.ClassName,
					YearAvailable:         max(class.Available, 1860),
					MaxLoadState:          5,
					ObjectName:            "Stairs (NE" + commaName + ")",
					PlatformConfiguration: rampConfiguration,
					UseCompanyColour:      true,
					HasFences:             true,
					OverrideOuter:         true,
					InnerPlatform:         inner,
					OuterPlatform:         outer,
					RailFenceSpriteID:     classSprites.SpriteMap["fence"],
					OuterPlatformSprite:   classSprites.SpriteMap["empty"],
				},
				{
				ID:                    baseObjectID + 24,
					BaseSpriteID:          classSprites.SpriteMap["empty"],
					ClassID:               class.ClassID,
					ClassName:             class.ClassName,
					YearAvailable:         max(class.Available, 1860),
					MaxLoadState:          5,
					ObjectName:            "Stairs (SW" + commaName + ")",
					PlatformConfiguration: rampConfiguration,
					UseCompanyColour:      true,
					HasFences:             true,
					OverrideOuter:         true,
					InnerPlatform:         inner,
					OuterPlatform:         outer,
					RailFenceSpriteID:     classSprites.SpriteMap["fence"],
					OuterPlatformSprite:   classSprites.SpriteMap["bare_stairs"],
				}}...)
			}

			if inner {
				thisClass = append(thisClass, nfo.Station{
					ID:                  baseObjectID + 21,
					BaseSpriteID:        classSprites.SpriteMap["bare_hut"],
					ClassID:             class.ClassID,
					ClassName:           class.ClassName,
					ObjectName:          "Hut" + bracketName,
					YearAvailable:       max(class.Available, 1890),
					MaxLoadState:        5,
					UseCompanyColour:    true,
					HasFences:           true,
					InnerPlatform:       inner,
					OuterPlatform:       outer,
					OverrideOuter:       true,
					OuterPlatformSprite: classSprites.SpriteMap["sign"],
					RailFenceSpriteID:   classSprites.SpriteMap["fence"],
				})
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
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
					},
					{
						ID:                baseObjectID + 20,
						BaseSpriteID:      classSprites.SpriteMap["bare_shelter_curved"],
						ClassID:           class.ClassID,
						ClassName:         class.ClassName,
						ObjectName:        "Shelter (curved" + commaName + ")",
						YearAvailable:     max(class.Available, 1926),
						MaxLoadState:      5,
						UseCompanyColour:  true,
						HasFences:         true,
						InnerPlatform:     inner,
						OuterPlatform:     outer,
						RailFenceSpriteID: classSprites.SpriteMap["fence"],
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
						RailFenceSpriteID: classSprites.SpriteMap["fence"],
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
						RailFenceSpriteID: classSprites.SpriteMap["fence"],
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
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
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
							RailFenceSpriteID:     classSprites.SpriteMap["fence"],
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
							RailFenceSpriteID:     classSprites.SpriteMap["fence"],
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
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
						PlatformConfiguration: rampConfiguration,
						AdditionalObjects: []nfo.AdditionalObject{
							{
								X:            5,
								Y:            4,
								Z:            12,
								SizeX:        5,
								SizeY:        8,
								SizeZ:        3,
								BaseSpriteID: footbridgeSprites[0].BaseSpriteID,
							},
						},
					}}...)
			}

			if class.ClassID != "TWF0" {
				footbridgeObjects := getFootbridgeBaseObject(footbridgeSprites[2].BaseSpriteID)

				if !inner || !outer {
					y := 0
					pillarSprite := 5
					if inner {
						y = 12
						pillarSprite = 6
					}

					footbridgeObjects = append(footbridgeObjects, getFootbridgeObject(y, footbridgeSprites[pillarSprite].BaseSpriteID))
				}

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
						InnerPlatform:         inner,
						OuterPlatform:         outer,
						PlatformConfiguration: rampConfiguration,
						AdditionalObjects:     footbridgeObjects,
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
					}}...)

				footbridgeObjects = getFootbridgeBaseObject(footbridgeSprites[1].BaseSpriteID)

				if !inner || !outer {
					y := 0
					pillarSprite := 3
					if inner {
						y = 12
						pillarSprite = 4
					}

					footbridgeObjects = append(footbridgeObjects, getFootbridgeObject(y, footbridgeSprites[pillarSprite].BaseSpriteID))
				}

				thisClass = append(thisClass, []nfo.Station{
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
						InnerPlatform:         inner,
						OuterPlatform:         outer,
						PlatformConfiguration: rampConfiguration,
						AdditionalObjects:     footbridgeObjects,
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
					},
				}...)
			}

			if i == 0 {

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
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
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
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
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
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
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
						RailFenceSpriteID:     classSprites.SpriteMap["fence"],
					},
				}...)
			}

			for _, station := range thisClass {
				thisStation := station
				outputObjects = append(outputObjects, &thisStation)
			}
		}

		// Wooden platforms do not have station halls
		if class.ClassID != "TWF0" {

			hall := nfo.StationHall{
				ID:                    class.BaseObjectID + (PLATFORM_TYPES * 3) + 1,
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

			outputObjects = append(outputObjects, &hall)
		}

		buffers := nfo.BufferStop{
			ID:               class.BaseObjectID + (PLATFORM_TYPES * 1) - 1,
			SpriteFilename:   fmt.Sprintf("%s_bufferstop", class.Filename),
			ClassID:          class.ClassID,
			ClassName:        class.ClassName,
			YearAvailable:    class.Available,
			BaseSpriteID:     classSprites.LastSpriteNumber + 14,
			ObjectName:       "Buffer Stop",
			UseCompanyColour: true,
		}

		outputObjects = append(outputObjects, &buffers)

		innerBuffers := nfo.BufferStop{
			ID:               class.BaseObjectID + (PLATFORM_TYPES * 2) - 1,
			SpriteFilename:   fmt.Sprintf("inner_%s_bufferstop", class.Filename),
			ClassID:          class.ClassID,
			ClassName:        class.ClassName,
			YearAvailable:    class.Available,
			BaseSpriteID:     classSprites.LastSpriteNumber + 14,
			ObjectName:       "Buffer Stop (inner)",
			UseCompanyColour: true,
			UseRailPresenceForSouth: true,
		}

		outputObjects = append(outputObjects, &innerBuffers)

		outerBuffers := nfo.BufferStop{
			ID:               class.BaseObjectID + (PLATFORM_TYPES * 3) - 1,
			SpriteFilename:   fmt.Sprintf("outer_%s_bufferstop", class.Filename),
			ClassID:          class.ClassID,
			ClassName:        class.ClassName,
			YearAvailable:    class.Available,
			BaseSpriteID:     classSprites.LastSpriteNumber + 14,
			ObjectName:       "Buffer Stop (outer)",
			UseCompanyColour: true,
			UseRailPresenceForNorth: true,
		}

		outputObjects = append(outputObjects, &outerBuffers)

		platforms := []nfo.FullTilePlatform{
			{
				SpriteFilename:   fmt.Sprintf("%s_concourse", class.Filename),
				ID:               class.BaseObjectID + (PLATFORM_TYPES * 3) + 5,
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				YearAvailable:    class.Available,
				ObjectName:       "Concourse",
				UseCompanyColour: true,
			},
			{
				SpriteFilename:   fmt.Sprintf("%s_concourse_shelter", class.Filename),
				ID:               class.BaseObjectID + (PLATFORM_TYPES * 3) + 6,
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				YearAvailable:    max(class.Available, 1860),
				ObjectName:       "Concourse with shelters",
				UseCompanyColour: true,
			},
		}

		for _, platform := range platforms {
			thisPlatform := platform
			outputObjects = append(outputObjects, &thisPlatform)
		}

		sort.Slice(outputObjects, func(i, j int) bool { return outputObjects[i].GetID() < outputObjects[j].GetID() })

		for _, object := range outputObjects {
			object.SetID(currentID)
			object.WriteToFile(&file)
			currentID = currentID + 1
		}

		outputObjects = make([]output_file.SortableFileWriter, 0)
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
			SpriteFilename:   "sir_william_tite",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Sir William Tite Station",
			YearAvailable:    1849,
			Width:            2,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "george_berkley",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Sir George Berkley Station",
			YearAvailable:    1853,
			Width:            2,
			UseCompanyColour: true,
		},
		{
			SpriteFilename:   "w_n_ashbee",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "W.N. Ashbee Station",
			YearAvailable:    1862,
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
			SpriteFilename:   "suburban",
			ClassID:          "TWFB",
			ClassName:        "Buildings",
			ObjectName:       "Suburban Station",
			YearAvailable:    1912,
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
			SpriteFilename: "straight_car_park",
			ClassID:        "TWFB",
			ClassName:      "Buildings",
			ObjectName:     "Car Park (straight)",
			YearAvailable:  1960,
			LoadStates:     3,
		},
		{
			SpriteFilename: "corner_car_park",
			ClassID:        "TWFB",
			ClassName:      "Buildings",
			ObjectName:     "Car Park (corner)",
			YearAvailable:  1960,
			LoadStates:     3,
		},
		{
			SpriteFilename: "entrance_car_park",
			ClassID:        "TWFB",
			ClassName:      "Buildings",
			ObjectName:     "Car Park (entrance)",
			YearAvailable:  1960,
			LoadStates:     3,
		},
		{
			SpriteFilename: "end_car_park",
			ClassID:        "TWFB",
			ClassName:      "Buildings",
			ObjectName:     "Car Park (end)",
			YearAvailable:  1960,
			LoadStates:     3,
		},
		{
			SpriteFilename: "shops",
			ClassID:        "TWFB",
			ClassName:      "Buildings",
			ObjectName:     "Shops",
			YearAvailable:  1970,
			LoadStates:     2,
		},
	}

	for _, building := range buildings {
		thisBuilding := building
		thisBuilding.ID = objectID
		outputObjects = append(outputObjects, &thisBuilding)
		objectID = objectID + 1
	}

	for _, building := range buildings {
		thisBuilding := building
		thisBuilding.ID = objectID
		thisBuilding.Reversed = true
		thisBuilding.ClassID = "TWFC"
		thisBuilding.ClassName = "Buildings (Reversed)"
		outputObjects = append(outputObjects, &thisBuilding)
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
			SpriteFilename:   "wp_arnside",
			ClassID:          "WAYP",
			ClassName:        "Waypoints",
			ObjectName:       "Arnside",
			YearAvailable:    1870,
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
			SpriteFilename:   "wp_shiplake",
			ClassID:          "WAYP",
			ClassName:        "Waypoints",
			ObjectName:       "Shiplake",
			YearAvailable:    1950,
			UseCompanyColour: true,
		},
	}

	for _, waypoint := range waypoints {
		thisWaypoint := waypoint
		thisWaypoint.ID = objectID
		outputObjects = append(outputObjects, &thisWaypoint)
		objectID = objectID + 1
	}

	sort.Slice(outputObjects, func(i, j int) bool { return outputObjects[i].GetID() < outputObjects[j].GetID() })

	for _, object := range outputObjects {
		object.SetID(currentID)
		object.WriteToFile(&file)
		currentID = currentID + 1
	}

	file.Output()
}

func getFootbridgeBaseObject(footbridgeSpriteID int) []nfo.AdditionalObject {
	return []nfo.AdditionalObject{
		{
			X:            6,
			Y:            2,
			Z:            13,
			SizeX:        5,
			SizeY:        8,
			SizeZ:        3,
			BaseSpriteID: footbridgeSpriteID,
		},
	}
}

func getFootbridgeObject(y int, footbridgeSpriteID int) nfo.AdditionalObject {
	return nfo.AdditionalObject{
		X:            6,
		Y:            y,
		Z:            0,
		SizeX:        3,
		SizeY:        3,
		SizeZ:        15,
		BaseSpriteID: footbridgeSpriteID,
	}
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
