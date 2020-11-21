package nfo

type Comment struct {
	Text string
}

func (c *Comment) GetComment() string {
	return c.Text
}

func (c *Comment) GetLines() []string {
	return []string{}
}
