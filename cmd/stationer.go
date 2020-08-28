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
	
	station1 := nfo.Station{
		SpriteFilename: "platforms_a",
		ClassID:        "TWF0",
		ClassName:      "British Stations",
		ObjectName:     "Platform A",
		ID: 0,
	}

	station1.WriteToFile(&file)

	station2 := nfo.Station{
		SpriteFilename: "platforms_sf",
		ClassID:        "TWF0",
		ClassName:      "British Stations",
		ObjectName:     "Simon Foster Platform",
		ID: 1,
	}

	station2.WriteToFile(&file)

	station3 := nfo.Station{
		SpriteFilename: "ramp_ne",
		ClassID:        "TWF0",
		ClassName:      "British Stations",
		ObjectName:     "Ramp (NE)",
		ID: 2,
	}

	station3.WriteToFile(&file)

	station4 := nfo.Station{
		SpriteFilename: "ramp_sw",
		ClassID:        "TWF0",
		ClassName:      "British Stations",
		ObjectName:     "Ramp (SW)",
		ID: 3,
	}

	station4.WriteToFile(&file)

	station5 := nfo.Station{
		SpriteFilename: "platforms_busy",
		ClassID:        "TWF0",
		ClassName:      "British Stations",
		ObjectName:     "Busy Platform (test)",
		ID: 4,
	}

	station5.WriteToFile(&file)

	station5 = nfo.Station{
		SpriteFilename: "platforms_modern",
		ClassID:        "TWF0",
		ClassName:      "British Stations",
		ObjectName:     "Modern Platform (test)",
		ID: 5,
	}

	station5.WriteToFile(&file)

	file.Output()
}
