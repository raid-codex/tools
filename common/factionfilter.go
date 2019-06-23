package common

type FactionFilter func(*Faction) bool

func FilterFactionSlug(slug string) FactionFilter {
	return func(faction *Faction) bool {
		return faction.Slug == slug
	}
}
