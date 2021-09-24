package book

type link struct {
	href string
	text string
	class string
}

func NewLink(href, text, class string) link {
	return link{href, text, class}
}

func (c link) Href() string {
	return c.href
}

func (c link) Text() string {
	return c.text
}

func (c link) Class() string {
	return c.class
}
