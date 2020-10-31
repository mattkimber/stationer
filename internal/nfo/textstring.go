package nfo

import (
	"fmt"
	bytes2 "github.com/mattkimber/stationer/internal/bytes"
)

const (
	TextStringTypeStationName = 0xC5
	TextStringTypeClassName   = 0xC4
)

type TextString struct {
	LanguageFile   int
	StationId      int
	TextStringType int
	Text           string
}

func (ts *TextString) GetComment() string {
	return "Text string definition"
}


func (ts *TextString) GetLines() []string {
	bytes := 7 + len(ts.Text)
	result := fmt.Sprintf("* %d 04 48 %s 01 %s %s \"%s\" 00",
		bytes,
		bytes2.GetByte(ts.LanguageFile),
		bytes2.GetByte(ts.StationId),
		bytes2.GetByte(ts.TextStringType),
		ts.Text,
	)
	return []string{result}
}
