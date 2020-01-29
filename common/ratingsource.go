package common

type AllRatings []*RatingSource

func (ar AllRatings) Compute() *Rating {
	if len(ar) == 0 {
		return &Rating{}
	} else if len(ar) == 1 {
		return ar[0].Rating
	}
	return &Rating{}
}

type RatingSource struct {
	Source string  `json:"source"`
	Rating *Rating `json:"rating"`
	Weight int     `json:"weight"`
}

func (rs *RatingSource) Sanitize() error { return rs.Rating.Sanitize() }
