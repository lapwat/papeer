package book

import (
	"fmt"

	"github.com/gosuri/uiprogress"
)

type progress struct {
	global      *uiprogress.Bar
	individuals []*uiprogress.Bar
}

func NewProgress(links []link) progress {
	uiprogress.Start()

	global := uiprogress.AddBar(len(links))
	global.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Chapters %d / %d", b.Current(), len(links))
	})

	individuals := []*uiprogress.Bar{}
	// hide individual bars if more than 50 chapters
	if len(links) <= 50 {
		for index, link := range links {
			bar := uiprogress.AddBar(1)
			barText := fmt.Sprintf("%d. %s", index+1, link.text)
			bar.AppendFunc(func(b *uiprogress.Bar) string {
				return barText
			})
			individuals = append(individuals, bar)
		}
	}

	return progress{global, individuals}
}

func (p *progress) IncrGlobal() {
	p.global.Incr()
}

func (p *progress) Incr(index int) {
	p.global.Incr()
	if len(p.individuals) > index {
		p.individuals[index].Incr()
	}
}
