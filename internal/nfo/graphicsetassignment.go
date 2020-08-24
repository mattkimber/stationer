package nfo

import "fmt"

type GraphicSetAssignment struct {
	IDs []int
	CargoSpecificSets []CargoToSet
	DefaultSet int
}

type CargoToSet struct {
	CargoType int
	Set int
}

func (gsa *GraphicSetAssignment) GetLines() []string {
	bytes := 5 + (len(gsa.CargoSpecificSets) * 3) + (len(gsa.IDs) * 2)

	ids := ""
	for _, id := range gsa.IDs {
		if len(ids) > 0 {
			ids += " "
		}
		ids += GetByte(id)
	}

	csets := ""
	for _, cset := range gsa.CargoSpecificSets {
		if len(csets) > 0 {
			csets += " "
		}
		csets += GetByte(cset.CargoType)
		csets += " "
		csets += GetShort(cset.Set)
	}

	result := fmt.Sprintf("* %d 03 04 %s %s %s %s %s",
		bytes,
		GetByte(len(gsa.IDs)),
		ids,
		GetByte(len(gsa.CargoSpecificSets)),
		csets,
		GetShort(gsa.DefaultSet),
	)

	return []string {result}
}