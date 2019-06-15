package templatefuncs

import (
	"fmt"
	"html/template"
	"strings"
)

var (
	FuncMap = template.FuncMap{
		"ReviewGrade":  reviewGrade,
		"ToLower":      strings.ToLower,
		"DisplayGrade": grade,
		"Percentage":   func(s float64) int64 { return int64(s * 100.0) },
		"TrustAsHtml":  func(s string) template.HTML { return template.HTML(s) },
		"dump":         func(v interface{}) string { return fmt.Sprintf("%+v", v) },
	}
)

var (
	ratingToRank = map[string]int{
		"SS": 5,
		"S":  4,
		"A":  3,
		"B":  2,
		"C":  1,
		"D":  0,
	}
	rarityToRank = map[string]int{
		"Legendary": 4,
		"Epic":      3,
		"Rare":      2,
		"Uncommon":  1,
		"Common":    0,
	}
)

func reviewGrade(gr float64) template.HTML {
	g := ""
	for v, r := range ratingToRank {
		if r == int(gr) {
			g = v
			break
		}
	}
	val := ``
	if g != "" {
		val = fmt.Sprintf(`<span class="champion-rating champion-rating-%s"><strong>%.1f</strong></span> `, g, gr)
	}
	return template.HTML(fmt.Sprintf("%s%s", val, grade(g)))
}

func grade(grade string) template.HTML {
	if grade == "" {
		return `<span class="champion-rating-none">No ranking yet</span>`
	}
	str := fmt.Sprintf(`<span class="champion-rating champion-rating-%s" title="%s">`, grade, gradeTitle(grade))
	for i := 0; i < 5; i++ {
		if i < ratingToRank[grade] {
			str += `<i class="fas fa-star"></i>`
		} else {
			str += `<i class="far fa-star"></i>`
		}
	}
	return template.HTML(str + `</span>`)
}

func gradeTitle(grade string) string {
	switch grade {
	case "D":
		return "not usable"
	case "C":
		return "viable"
	case "B":
		return "good"
	case "A":
		return "exceptional"
	case "S":
		return "top tier"
	case "SS":
		return "god tier"
	}
	return ""
}
