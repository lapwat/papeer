package book

type link struct {
	href string
	text string
}

func NewLink(href, text string) link {
	return link{href, text}
}

func (c link) Href() string {
	return c.href
}

func (c link) Text() string {
	return c.text
}
