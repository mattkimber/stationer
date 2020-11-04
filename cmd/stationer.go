package main

import (
	"flag"
	"fmt"
	"github.com/mattkimber/stationer/internal/nfo"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

func init() {
	flag.Parse()
}

type StationClass struct {
	Filename  string
	ClassID   string
	ClassName string
	Available int
}

func main() {
	file := nfo.File{}
	file.AddElement(&nfo.Header{
		Initials:    "TWF",
		SetID:       8,
		SetName:     "Timberwolf's Stations 0.2.1 (alpha)",
		Description: "A set of British-style railway stations",
		Version:     1,
		MinVersion:  1,
	})

	file.AddElement(&nfo.CargoTypeTable{Cargos: []string{"PASS", "MAIL"}})

	classes := []StationClass{
		{Filename: "wooden", ClassID: "TWF0", ClassName: "Wooden Platforms", Available: 0},
		{Filename: "concrete", ClassID: "TWF1", ClassName: "Concrete Platforms", Available: 1860},
		{Filename: "modern", ClassID: "TWF2", ClassName: "Modern Platforms", Available: 1970},
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

	stations := make([]nfo.Station, 0)

	for _, class := range classes {

		thisClass := []nfo.Station{
			{
				SpriteFilename:   class.Filename + "_empty",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform",
				YearAvailable:    class.Available,
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:   class.Filename + "_empty",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform (inner)",
				YearAvailable:    class.Available,
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    false,
			},
			{
				SpriteFilename:   class.Filename + "_empty",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform (outer)",
				YearAvailable:    class.Available,
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    false,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:   class.Filename + "_sign",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform with sign",
				YearAvailable:    max(class.Available, 1845),
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:   class.Filename + "_sign",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform with sign (inner)",
				YearAvailable:    max(class.Available, 1845),
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    false,
			},
			{
				SpriteFilename:   class.Filename + "_sign",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform with sign (outer)",
				YearAvailable:    max(class.Available, 1845),
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    false,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:   class.Filename + "_benches",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform with benches",
				YearAvailable:    max(class.Available, 1840),
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:   class.Filename + "_benches",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform with benches (inner)",
				YearAvailable:    max(class.Available, 1840),
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    false,
			},
			{
				SpriteFilename:   class.Filename + "_benches",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Platform with benches (outer)",
				YearAvailable:    max(class.Available, 1840),
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    false,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:   class.Filename + "_bare_shelter_traditional",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Shelter (traditional)",
				YearAvailable:    max(class.Available, 1860),
				MaxLoadState:     5,
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:   class.Filename + "_bare_shelter_traditional",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Shelter (traditional, inner)",
				YearAvailable:    max(class.Available, 1860),
				MaxLoadState:     5,
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    true,
				OuterPlatform:    false,
			},
			{
				SpriteFilename:   class.Filename + "_bare_shelter_traditional",
				ClassID:          class.ClassID,
				ClassName:        class.ClassName,
				ObjectName:       "Shelter (traditional, outer)",
				YearAvailable:    max(class.Available, 1860),
				MaxLoadState:     5,
				UseCompanyColour: true,
				HasFences:        true,
				InnerPlatform:    false,
				OuterPlatform:    true,
			},
			{
				SpriteFilename:        class.Filename + "_ramp_ne",
				ClassID:               class.ClassID,
				ClassName:             class.ClassName,
				YearAvailable:         class.Available,
				MaxLoadState:          5,
				InnerPlatform:         true,
				OuterPlatform:         true,
				ObjectName:            "Ramp (NE)",
				PlatformConfiguration: rampConfiguration,
				UseCompanyColour:      true,
				HasFences:             true,
			},
			{
				SpriteFilename:        class.Filename + "_ramp_ne",
				ClassID:               class.ClassID,
				ClassName:             class.ClassName,
				YearAvailable:         class.Available,
				MaxLoadState:          5,
				InnerPlatform:         true,
				OuterPlatform:         false,
				ObjectName:            "Ramp (NE, inner)",
				PlatformConfiguration: rampConfiguration,
				UseCompanyColour:      true,
				HasFences:             true,
			},
			{
				SpriteFilename:        class.Filename + "_ramp_ne",
				ClassID:               class.ClassID,
				ClassName:             class.ClassName,
				YearAvailable:         class.Available,
				MaxLoadState:          5,
				InnerPlatform:         false,
				OuterPlatform:         true,
				ObjectName:            "Ramp (NE, outer)",
				PlatformConfiguration: rampConfiguration,
				UseCompanyColour:      true,
				HasFences:             true,
			},
			{
				SpriteFilename:        class.Filename + "_ramp_sw",
				ClassID:               class.ClassID,
				ClassName:             class.ClassName,
				YearAvailable:         class.Available,
				MaxLoadState:          5,
				InnerPlatform:         true,
				OuterPlatform:         true,
				ObjectName:            "Ramp (SW)",
				PlatformConfiguration: rampConfiguration,
				UseCompanyColour:      true,
				HasFences:             true,
			},
			{
				SpriteFilename:        class.Filename + "_ramp_sw",
				ClassID:               class.ClassID,
				ClassName:             class.ClassName,
				YearAvailable:         class.Available,
				MaxLoadState:          5,
				InnerPlatform:         true,
				OuterPlatform:         false,
				ObjectName:            "Ramp (SW, inner)",
				PlatformConfiguration: rampConfiguration,
				UseCompanyColour:      true,
				HasFences:             true,
			},
			{
				SpriteFilename:        class.Filename + "_ramp_sw",
				ClassID:               class.ClassID,
				ClassName:             class.ClassName,
				YearAvailable:         class.Available,
				MaxLoadState:          5,
				InnerPlatform:         false,
				OuterPlatform:         true,
				ObjectName:            "Ramp (SW, outer)",
				PlatformConfiguration: rampConfiguration,
				UseCompanyColour:      true,
				HasFences:             true,
			},
			{
				SpriteFilename:        class.Filename + "_bare_footbridge",
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
				AdditionalObjects: []nfo.AdditionalObject{{
					X:              5,
					Y:              4,
					Z:              13,
					SizeX:          5,
					SizeY:          8,
					SizeZ:          3,
					SpriteFilename: "footbridge",
				}},
			},
		}
		stations = append(stations, thisClass...)
	}

	objectID := 0

	for _, station := range stations {
		station.ID = objectID
		station.WriteToFile(&file)
		objectID = objectID + 1
	}

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
