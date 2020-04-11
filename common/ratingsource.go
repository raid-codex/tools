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
	gTotal := 0
	gDivideBy := 0
	for i := 0; i < indV.NumField(); i++ {
		if tag := indV.Type().Field(i).Tag.Get("json"); tag == "overall" {
			continue
		}
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
		gTotal += total
		gDivideBy += divideBy
	}
	if gDivideBy > 0 {
		rating.Overall = overallRatioToRank(float32(gTotal) / float32(gDivideBy))
	}
	return rating
}

type RatingSource struct {
	Source string  `json:"source"`
	Rating *Rating `json:"rating"`
	Weight int     `json:"weight"`
}

func (rs *RatingSource) Sanitize() error {
	switch rs.Source {
	case "ayumilove":
		rs.Weight = 2
	}
	return rs.Rating.Sanitize()
}
