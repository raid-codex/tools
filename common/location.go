package common

type Location struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (l *Location) Sanitize() error {
	return nil
}
