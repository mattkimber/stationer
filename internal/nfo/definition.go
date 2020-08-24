package nfo

import "fmt"

type Definition struct {
	StationID int
	properties []Property
}

func (d *Definition) AddProperty(property Property) {
	if d.properties == nil {
		d.properties = make([]Property, 0)
	}

	d.properties = append(d.properties, property)
}

func (d *Definition) GetLines() []string {
	bytes := 5
	output := ""

	for _, p := range d.properties {
		bytes += p.GetBytes()
		output += p.GetString()
	}

	result := fmt.Sprintf("* %d 00 04 %s 01 %s %s",
		bytes,
		GetByte(len(d.properties)),
		GetByte(d.StationID),
		output)

	return []string { result }
}
