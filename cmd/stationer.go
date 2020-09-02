package main

import (
	"flag"
	"github.com/mattkimber/stationer/internal/nfo"
	"github.com/mattkimber/stationer/internal/nfo/properties"
)

func init() {
	flag.Parse()
}

type StationClass struct {
	Filename string
	ClassID string
	ClassName string
}

func main() {
	file := nfo.File{}
	file.AddElement(&nfo.Header{
		Initials:    "TWF",
		SetID:       8,
		SetName:     "Timberwolf's Stations",
		Description: "A set of British-style railway stations",
		Version: 1,
		MinVersion: 1,
	})

	file.AddElement(&nfo.CargoTypeTable{Cargos: []string{"PASS", "MAIL"}})

	classes := []StationClass{
		{ Filename: "concrete", ClassID: "TWF0", ClassName: "Concrete Platforms" },
		{ Filename: "modern", ClassID: "TWF1", ClassName: "Modern Platforms" },
	}

	stations := make([]nfo.Station, 0)

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
		Lengths:   properties.PlatformBitmask{
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
		thisClass := []nfo.Station{
			{
				SpriteFilename: class.Filename + "_empty",
				ClassID:        class.ClassID,
				ClassName:      class.ClassName,
				ObjectName:     "Platform",
			},
			{
				SpriteFilename: class.Filename + "_sign",
				ClassID:        class.ClassID,
				ClassName:      class.ClassName,
				ObjectName:     "Platform with sign",
				UseCompanyColour: true,
			},
			{
				SpriteFilename: class.Filename + "_benches",
				ClassID:        class.ClassID,
				ClassName:      class.ClassName,
				ObjectName:     "Platform with benches",
			},
			{
				SpriteFilename: class.Filename + "_ramp_ne",
				ClassID:        class.ClassID,
				ClassName:      class.ClassName,
				ObjectName:     "Ramp (NE)",
				PlatformConfiguration: rampConfiguration,
			},
			{
				SpriteFilename: class.Filename + "_ramp_sw",
				ClassID:        class.ClassID,
				ClassName:      class.ClassName,
				ObjectName:     "Ramp (SW)",
				PlatformConfiguration: rampConfiguration,
			},
		}

		stations = append(stations, thisClass...)
	}


	for idx, station := range stations {
		station.ID = idx
		station.WriteToFile(&file)
	}

	file.Output()
}
