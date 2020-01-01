package common

import "time"

type Video struct {
	Source    string `json:"source"`
	ID        string `json:"id"`
	Author    string `json:"author"`
	DateAdded string `json:"date_added"`
}

func (v *Video) Sanitize() error {
	if v.DateAdded == "" {
		v.DateAdded = time.Now().Format(time.RFC3339)
	}
	return nil
}
