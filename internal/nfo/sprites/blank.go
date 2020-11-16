package sprites

type Blank struct {
	Size int
}

func (b *Blank) GetComment() string {
	return "Blank pseudosprite"
}

func (b *Blank) GetLines() []string {
	result := make([]string, b.Size)

	for i := 0; i < b.Size; i++ {
		// Sprite that will never be used
		result[i] = "-1 * 1 00"
	}

	return result
}
