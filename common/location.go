package common

type Location struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (l *Location) Sanitize() error {
	return nil
}

func ConvertLocation(s string) string {
	if v, ok := map[string]string{
		"dungeons": "dungeon",
	}[s]; ok {
		return v
	}
	return s
}
