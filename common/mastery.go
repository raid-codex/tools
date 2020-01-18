package common

import "strings"

type Mastery struct {
	Name        string `json:"name"`
	Star        int64  `json:"star"`
	Description string `json:"description"`
	Tree        uint8  `json:"tree"`
	Level       uint8  `json:"level"`
	ScrollType  uint8  `json:"scroll_type"`
	Unlock      uint64 `json:"unlock"`
	ImageSlug   string `json:"image_slug"`
	Slug        string `json:"slug"`
}

type MasteryList []*Mastery

func (ml MasteryList) Sort() {
}

type MasteryFilter func(*Mastery) bool

func FilterMasteryLowercasedName(name string) MasteryFilter {
	name = strings.ToLower(name)
	return func(m *Mastery) bool {
		return strings.ToLower(m.Name) == name
	}
}
