package book

import "time"

type link struct {
	Href string     `json:"url"`
	Text string     `json:"name"`
	Date *time.Time `json:"date"`
}

func NewLink(href, text string, date *time.Time) link {
	return link{href, text, date}
}
