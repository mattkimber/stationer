package properties

import "fmt"

// CLass ID
type ClassID struct {
	ID string
}

func (c *ClassID) GetBytes() int {
	return 1 + len(c.ID)
}

func (c *ClassID) GetString() string {
	return fmt.Sprintf("08 \"%s\" ", c.ID)
}
