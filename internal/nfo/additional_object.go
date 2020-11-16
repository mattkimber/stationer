package nfo

type AdditionalObject struct {
	X               int
	Y               int
	Z               int
	SizeX           int
	SizeY           int
	SizeZ           int
	BaseSpriteID    int
	InvertDirection bool
}

func (ao *AdditionalObject) GetBaseSpriteNumber(s *Station) int {
	if s.UseCompanyColour {
		return COMPANY_COLOUR_SPRITE + ao.BaseSpriteID
	}

	return CUSTOM_SPRITE + ao.BaseSpriteID
}
