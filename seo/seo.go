package seo

import (
	"encoding/json"
)

type SEO struct {
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Keywords       []string          `json:"keywords"`
	StructuredData []json.RawMessage `json:"structured_data"`
}
