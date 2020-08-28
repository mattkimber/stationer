package nfo

import (
	"fmt"
	"log"
)

const (
	HEADER_ACTION = 8
	NFO_VERSION = 8
	ACTION_LENGTH = 1
	VERSION_LENGTH = 1
	GRFID_LENGTH = 4
	NULL_LENGTH = 1
)

type Header struct {
	Initials string
	SetID int
	SetName string
	Description string
	Version int
	MinVersion int
}

func (h *Header) GetLines() []string {
	if len(h.Initials) != 3 {
		log.Fatalf("Initials must be exactly 3 characters")
	}

	if h.SetID > 255 {
		log.Fatalf("Set ID can be at most 255")
	}

	bytes := ACTION_LENGTH + VERSION_LENGTH + GRFID_LENGTH + len(h.SetName) + NULL_LENGTH + len(h.Description) + NULL_LENGTH

	action14 := fmt.Sprintf("* 0 14 \"C\" \"INFO\" \"B\" \"PALS\" 01 00 \"D\" \"B\" \"VRSN\" 04 00 %s \"B\" \"MINV\" 04 00 %s 00 00", GetWord(h.Version), GetWord(h.MinVersion))

	action8 := fmt.Sprintf("* %d %s %s \"%s\" %s \"%s\" %s \"%s\" %s",
		bytes,
		GetByte(HEADER_ACTION),
		GetByte(NFO_VERSION),
		h.Initials,
		GetByte(h.SetID),
		h.SetName,
		GetByte(0),
		h.Description,
		GetByte(0),
	)
	return []string{ action14, action8 }
}