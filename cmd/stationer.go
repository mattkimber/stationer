package main

import (
	"flag"
	"github.com/mattkimber/stationer/internal/nfo"
)

func init() {
	flag.Parse()
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

	stations := []nfo.Station{
		{
			SpriteFilename: "concrete_empty",
			ClassID:        "TWF0",
			ClassName:      "Concrete Platforms",
			ObjectName:     "Platform",
		},
		{
			SpriteFilename: "concrete_sign",
			ClassID:        "TWF0",
			ClassName:      "Concrete Platforms",
			ObjectName:     "Platform with sign",
		},
		{
			SpriteFilename: "concrete_benches",
			ClassID:        "TWF0",
			ClassName:      "Concrete Platforms",
			ObjectName:     "Platform with benches",
		},
		{
			SpriteFilename: "modern_empty",
			ClassID:        "TWF1",
			ClassName:      "Modern Platforms",
			ObjectName:     "Platform",
		},
		{
			SpriteFilename: "modern_sign",
			ClassID:        "TWF1",
			ClassName:      "Modern Platforms",
			ObjectName:     "Platform with sign",
		},
		{
			SpriteFilename: "modern_benches",
			ClassID:        "TWF1",
			ClassName:      "Modern Platforms",
			ObjectName:     "Platform with benches",
		},
	}

	for idx, station := range stations {
		station.ID = idx
		station.WriteToFile(&file)
	}

	file.Output()
}
