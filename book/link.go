package book

type link struct {
	Href string `json:"url"`
	Text string `json:"name"`
}

func NewLink(href, text string) link {
	return link{href, text}
}
