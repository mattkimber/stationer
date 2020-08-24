package nfo

import "fmt"

const (
	TextStringTypeStationName = 0xC5
	TextStringTypeClassName = 0xC4
)

type TextString struct {
	LanguageFile int
	StationId int
	TextStringType int
	Text string
}

func (ts *TextString) GetLines() []string {
	bytes := 7 + len(ts.Text)
	result := fmt.Sprintf("* %d 04 48 %s 01 %s %s \"%s\" 00",
		bytes,
		GetByte(ts.LanguageFile),
		GetByte(ts.StationId),
		GetByte(ts.TextStringType),
		ts.Text,
		)
	return []string { result }
}
