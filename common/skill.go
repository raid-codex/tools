package common

type Skill struct {
	Name           string `json:"name"`
	RawDescription string `json:"raw_description"`
	Slug           string `json:"slug"`
}

func (s *Skill) Sanitize() error {
	s.Slug = GetLinkNameFromSanitizedName(s.Name)
	return nil
}
