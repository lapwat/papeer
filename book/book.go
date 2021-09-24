package book

type book struct {
	name     string
	author   string
	chapters []chapter
}

func New(name, author string) book {
	return book{name, author, []chapter{}}
}

func (b *book) AddChapter(c chapter) {
	b.chapters = append(b.chapters, c)
}

func (b book) Name() string {
	return b.name
}

func (b book) Author() string {
	return b.name
}

func (b *book) Chapters() []chapter {
	return b.chapters
}
