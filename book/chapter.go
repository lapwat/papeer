package book

type chapter struct {
	name    string
	author  string
	content string
}

func NewChapter(name, author, content string) chapter {
	return chapter{name, author, content}
}

func (c chapter) Name() string {
	return c.name
}

func (c chapter) Author() string {
	return c.author
}

func (c chapter) Content() string {
	return c.content
}
