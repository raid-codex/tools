package common

import "reflect"

type AllRatings []*RatingSource

func (ar AllRatings) Compute() *Rating {
	if len(ar) == 0 {
		return &Rating{}
	} else if len(ar) == 1 {
		return ar[0].Rating
	}
	rating := &Rating{}
	v := reflect.ValueOf(rating)
	indV := reflect.Indirect(v)
	for i := 0; i < indV.NumField(); i++ {
		divideBy := 0
		total := 0
		for _, r := range ar {
			value := reflect.Indirect(reflect.ValueOf(r.Rating)).Field(i).String()
			if _, ok := rankToInt[value]; !ok {
				continue
			}
			total += r.Weight * rankToInt[value]
			divideBy += r.Weight
		}
		if divideBy > 0 {
			indV.Field(i).SetString(intToRank[int(float32(total)/float32(divideBy))])
		}
	}
	return rating
}

type RatingSource struct {
	Source string  `json:"source"`
	Rating *Rating `json:"rating"`
	Weight int     `json:"weight"`
}

func (rs *RatingSource) Sanitize() error {
	return rs.Rating.Sanitize()
}
