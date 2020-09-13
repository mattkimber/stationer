package nfo

import (
	"fmt"
	bytes2 "github.com/mattkimber/stationer/internal/bytes"
)

type GraphicSetAssignment struct {
	IDs               []int
	CargoSpecificSets []CargoToSet
	DefaultSet        int
}

type CargoToSet struct {
	CargoType int
	Set       int
}

func (gsa *GraphicSetAssignment) GetLines() []string {
	bytes := 5 + (len(gsa.CargoSpecificSets) * 3) + (len(gsa.IDs) * 2)

	ids := ""
	for _, id := range gsa.IDs {
		if len(ids) > 0 {
			ids += " "
		}
		ids += bytes2.GetByte(id)
	}

	csets := ""
	for _, cset := range gsa.CargoSpecificSets {
		if len(csets) > 0 {
			csets += " "
		}
		csets += bytes2.GetByte(cset.CargoType)
		csets += " "
		csets += bytes2.GetWord(cset.Set)
	}

	result := fmt.Sprintf("* %d 03 04 %s %s %s %s %s",
		bytes,
		bytes2.GetByte(len(gsa.IDs)),
		ids,
		bytes2.GetByte(len(gsa.CargoSpecificSets)),
		csets,
		bytes2.GetWord(gsa.DefaultSet),
	)

	return []string{result}
}
