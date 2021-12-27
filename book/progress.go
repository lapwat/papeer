package book

import (
	"fmt"
	"strings"

	"github.com/gosuri/uiprogress"
)

type progress struct {
	global      *uiprogress.Bar
	individuals []*uiprogress.Bar
}

func NewProgress(links []link, parent string, depth int) progress {
	uiprogress.Start()

	global := uiprogress.AddBar(len(links))
	indentGlobal := strings.Repeat("> ", depth)
	global.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%v%v (%v / %v)", indentGlobal, parent, b.Current(), len(links))
	})

	// hide individual bars if more than 50 chapters
	individuals := []*uiprogress.Bar{}
	indent := strings.Repeat("- ", depth)
	if len(links) <= 50 {
		for index, link := range links {
			bar := uiprogress.AddBar(1)
			barText := fmt.Sprintf("%v#%v %v", indent, index+1, link.Text())
			bar.AppendFunc(func(b *uiprogress.Bar) string {
				return barText
			})
			individuals = append(individuals, bar)
		}
	}

	return progress{global, individuals}
}

func (p *progress) IncrementGlobal() {
	p.global.Incr()
}

func (p *progress) Increment(index int) {
	p.IncrementGlobal()
	if len(p.individuals) > index {
		p.individuals[index].Incr()
	}
}

func (p *progress) UpdateName(index int, name string) {
	if len(p.individuals) > index {
		barText := fmt.Sprintf("%s", name)
		p.individuals[index].AppendFunc(func(b *uiprogress.Bar) string {
			return barText
		})
	}
}
