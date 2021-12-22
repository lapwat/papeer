package book

type chapter struct {
	body        string
	name        string
	author      string
	content     string
	subChapters []chapter
	config      *ScrapeConfig
}

func NewChapter(body, name, author, content string, subChapters []chapter, config *ScrapeConfig) chapter {
	return chapter{body, name, author, content, subChapters, config}
}

func (c chapter) Body() string {
	return c.body
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

func (c chapter) SubChapters() []chapter {
	return c.subChapters
}
