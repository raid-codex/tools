package common

type FusionFilter func(*Fusion) bool

func FilterFusionSlug(slug string) FusionFilter {
	return func(fusion *Fusion) bool { return fusion.Slug == slug }
}
