package common

type Video struct {
	Source string `json:"source"`
	ID     string `json:"id"`
	Author string `json:"author"`
}

func (v *Video) Sanitize() error {
	return nil
}
