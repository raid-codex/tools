package templatefuncs

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
)

var (
	rootUrl = "https://raid-codex.com"
	FuncMap = template.FuncMap{
		"ReviewGrade":  reviewGrade,
		"ToLower":      strings.ToLower,
		"DisplayGrade": grade,
		"Percentage":   func(s float64) int64 { return int64(s * 100.0) },
		"TrustAsHtml":  func(s string) template.HTML { return template.HTML(s) },
		"dump":         func(v interface{}) string { return fmt.Sprintf("%+v", v) },
		"synergy_raw_title": func(s common.SynergyContextKey) string {
			if v, ok := common.SynergyData[s]; ok {
				return v.Title
			}
			panic(fmt.Errorf("synergy Title not found: %s", s))
		},
		"synergy_raw_description": func(s common.SynergyContextKey) string {
			if v, ok := common.SynergyData[s]; ok {
				return v.RawDescription
			}
			panic(fmt.Errorf("synergy RawDescription not found: %s", s))
		},
		"getChampions": func(s []string) common.ChampionList {
			champions, errChampions := common.GetChampions(func(champion *common.Champion) bool {
				for _, c := range s {
					if c == champion.Slug {
						return true
					}
				}
				return false
			})
			if errChampions != nil {
				panic(errChampions)
			}
			return champions
		},
		"championImage": func(slug string) string {
			return fmt.Sprintf("%s/wp-content/uploads/champions/image-champion-%s.jpg", rootUrl, slug)
		},
		"championThumbnail": func(slug string) string {
			return fmt.Sprintf("%s/wp-content/uploads/champion-thumbnails/image-champion-small-%s.jpg", rootUrl, slug)
		},
		"championThumbnailFallback": func(slug string) string {
			champions, _ := common.GetChampions(func(champion *common.Champion) bool {
				return champion.Slug == slug
			})
			if len(champions) != 1 {
				panic(champions)
			}
			img, err := utils.ImageFallback(
				fmt.Sprintf("%s/wp-content/uploads/hashed-img/%s.png", rootUrl, champions[0].Thumbnail),
				fmt.Sprintf("%s/wp-content/uploads/champion-thumbnails/image-champion-small-%s.jpg", rootUrl, slug),
				blankImage,
			)
			if err != nil {
				panic(err)
			}
			return img
		},
		"websiteLink": func(websiteLink string) string {
			return fmt.Sprintf("%s%s", rootUrl, websiteLink)
		},
		"championImageFallback": func(slug string) string {
			champions, _ := common.GetChampions(func(champion *common.Champion) bool {
				return champion.Slug == slug
			})
			if len(champions) != 1 {
				panic(fmt.Sprintf("found champions: %+v for slug %s", champions, slug))
			}
			img, err := utils.ImageFallback(
				fmt.Sprintf("%s/wp-content/uploads/champions/image-champion-%s.jpg", rootUrl, slug),
				fmt.Sprintf("%s/wp-content/uploads/hashed-img/%s.png", rootUrl, champions[0].Thumbnail),
				fmt.Sprintf("%s/wp-content/uploads/champion-thumbnails/image-champion-small-%s.jpg", rootUrl, slug),
				blankImage,
			)
			if err != nil {
				panic(err)
			}
			return img
		},
		"skillImageFallback": func(slug string) string {
			img, err := utils.ImageFallback(
				fmt.Sprintf("%s/wp-content/uploads/hashed-img/%s.png", rootUrl, slug),
				blankImage,
			)
			if err != nil {
				panic(err)
			}
			return img
		},
		"safeAttr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
		"effectImage": func(se *common.StatusEffect) template.HTML {
			img, err := utils.ImageFallback(
				fmt.Sprintf("%s/wp-content/uploads/status-effects/%s.png", rootUrl, se.ImageSlug),
				blankImage,
			)
			if err != nil {
				panic(err)
			}
			return template.HTML(fmt.Sprintf(
				`<img src="%s" title="%s" alt="%s">`, img, se.RawDescription, se.Type,
			))
		},
		"factionImage": func(faction *common.Faction) template.HTML {
			img, err := utils.ImageFallback(
				fmt.Sprintf("%s/wp-content/uploads/factions/%s.png", rootUrl, faction.ImageSlug),
				blankImage,
			)
			if err != nil {
				panic(err)
			}
			return template.HTML(fmt.Sprintf(
				`<img src="%s" title="%s" alt="%s">`, img, faction.Name, faction.Name,
			))
		},
		"getHitsOfSkill": func(skill *common.Skill) int64 {
			if len(skill.Upgrades) == 0 {
				return 0
			}
			return skill.Upgrades[len(skill.Upgrades)-1].Hits
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"add": func(a, b int) int {
			return a + b
		},
		"repeat": func(s int64) []int64 {
			a := make([]int64, s)
			for i := int64(0); i < s; i++ {
				a[i] = i + 1
			}
			return a
		},
		"getRootFusion": func(fusions map[string]*common.Fusion, fusionSlug string) *common.Fusion {
			fusion := fusions[fusionSlug]
			for fusion.ParentFusionSlug != nil {
				fusion = fusions[*fusion.ParentFusionSlug]
			}
			return fusion
		},
		"displayLocations": func(locations []string) []string {
			v := map[string]string{
				"clan-boss": "Clan Boss",
				"dungeon":   "Dungeon",
				"arena":     "Arena",
				"campaign":  "Campaign",
			}
			l := []string{}
			for _, location := range locations {
				l = append(l, v[location])
			}
			return l
		},
		"displaySet": func(sets []string) []string {
			s := []string{}
			for _, set := range sets {
				s = append(s, strings.Title(strings.Replace(set, "-", " ", -1)))
			}
			return s
		},
		"joinStrings": func(s []string, sep string) string {
			return strings.Join(s, sep)
		},
	}
)

const (
	blankImage = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="
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
	if gr == 0.0 {
		return template.HTML(`<span class="champion-rating-none">No ranking yet</span>`)
	}
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
		return template.HTML(`<span class="champion-rating-none">No ranking yet</span>`)
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
