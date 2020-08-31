package nfo

import "fmt"

type CargoTypeTable struct {

	Cargos []string
}


func (ct *CargoTypeTable) GetLines() []string {
	bytes := 6
	output := ""

	for _, cargo := range ct.Cargos {
		output = output + "\"" + cargo + "\" "
		bytes += 4
	}

	result := fmt.Sprintf("* %d 00 08 %s %s 00 09 %s",
		bytes,
		GetByte(1), // Changing 1 property
		GetByte(len(ct.Cargos)), // n cargo types
		output)

	return []string { result }
}
