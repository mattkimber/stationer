package sprites

import "fmt"

type Blank struct {
	Size int
	Name string
}

func (b *Blank) GetComment() string {
	return fmt.Sprintf("Blank pseudosprite (%s)", b.Name)
}

func (b *Blank) GetLines() []string {
	result := make([]string, b.Size)

	for i := 0; i < b.Size; i++ {
		// Sprite that will never be used
		result[i] = "-1 * 1 00"
	}

	return result
}
