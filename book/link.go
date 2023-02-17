package book

type link struct {
	Href string `json:"href"`
	Text string `json:"name"`
}

func NewLink(href, text string) link {
	return link{href, text}
}
