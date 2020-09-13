package nfo

import (
	"fmt"
	bytes2 "github.com/mattkimber/stationer/internal/bytes"
	properties2 "github.com/mattkimber/stationer/internal/nfo/properties"
)

type Definition struct {
	StationID  int
	properties []properties2.Property
}

func (d *Definition) AddProperty(property properties2.Property) {
	if d.properties == nil {
		d.properties = make([]properties2.Property, 0)
	}

	d.properties = append(d.properties, property)
}

func (d *Definition) GetLines() []string {
	bytes := 5
	output := ""

	for _, p := range d.properties {
		bytes += p.GetBytes()

		// Make the output a little more readable by adding a new line for each property
		output += "\n    " + p.GetString()
	}

	result := fmt.Sprintf("* %d 00 04 %s 01 %s %s",
		bytes,
		bytes2.GetByte(len(d.properties)),
		bytes2.GetByte(d.StationID),
		output)

	return []string{result}
}
