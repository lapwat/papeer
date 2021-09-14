package book

type Book struct {
	name     string
	author   string
	chapters []chapter
}

func New(name, author string) Book {
	return Book{name, author, []chapter{}}
}

func (b *Book) AddChapter(c chapter) {
	b.chapters = append(b.chapters, c)
}

func (b Book) Name() string {
	return b.name
}

func (b Book) Author() string {
	return b.name
}

func (b *Book) Chapters() []chapter {
	return b.chapters
}
