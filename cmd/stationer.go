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
	
	file.AddElement(&nfo.Sprites{
		{
			Filename: "platforms_a_8bpp.png",
			X:        0,
			Y:        0,
			Width:    64,
			Height:   55,
			XRel:     -64/2,
			YRel:     -55/2,
		},
		{
			Filename: "platforms_b_8bpp.png",
			X:        0,
			Y:        0,
			Width:    64,
			Height:   55,
			XRel:     -64/2,
			YRel:     -55/2,
		},
		{
			Filename: "platforms_a_8bpp.png",
			X:        72,
			Y:        0,
			Width:    64,
			Height:   55,
			XRel:     -64/2,
			YRel:     -55/2,
		},
		{
			Filename: "platforms_b_8bpp.png",
			X:        72,
			Y:        0,
			Width:    64,
			Height:   55,
			XRel:     -64/2,
			YRel:     -55/2,
		},
	})

	def := &nfo.Definition{StationID: 0}
	def.AddProperty(&nfo.ClassID{ID: "TWF0"})
	def.AddProperty(&nfo.SpriteLayout{
		EastWest:   nfo.SpriteDirection{
			GroundSprite: 1012,
			Foreground:   nfo.BoundingBox{X: 16, Y: 5, Z: 5},
			Background:   nfo.BoundingBox{X: 16, Y: 5, Z: 3},
		},
		NorthSouth: nfo.SpriteDirection{
			GroundSprite: 1011,
			Foreground:   nfo.BoundingBox{X: 16, Y: 5, Z: 5},
			Background:   nfo.BoundingBox{X: 16, Y: 5, Z: 3},
		},
	})

	file.AddElement(def)
	
	file.AddElement(&nfo.StationSet{
		SetID:         0,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})

	file.AddElement(&nfo.StationSet{
		SetID:         1,
		NumLittleSets: 0,
		NumLotsSets:   1,
		SpriteSets:    []int{0},
	})
	
	file.AddElement(&nfo.GraphicSetAssignment{
		IDs:               []int {0},
		CargoSpecificSets: []nfo.CargoToSet{{
			CargoType: 254,
			Set:       1,
		}},
		DefaultSet:        0,
	})

	file.AddElement(&nfo.TextString{
		LanguageFile:   255,
		StationId:      0,
		TextStringType: nfo.TextStringTypeStationName,
		Text:           "Station Tile",
	})

	file.AddElement(&nfo.TextString{
		LanguageFile:   255,
		StationId:      0,
		TextStringType: nfo.TextStringTypeClassName,
		Text:           "Timberwolf",
	})

	file.Output()
}
