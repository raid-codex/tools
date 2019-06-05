package common

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type StatusEffect struct {
	EffectType     string  `json:"effect_type"`
	Type           string  `json:"type"`
	Value          float64 `json:"value"`
	ImageSlug      string  `json:"image_slug"`
	Slug           string  `json:"slug"`
	WebsiteLink    string  `json:"website_link"`
	Extra          bool    `json:"extra"`
	RawDescription string  `json:"raw_description"`
}

func (se *StatusEffect) Sanitize() error {
	if se.Slug == "" {
		se.Slug = GetLinkNameFromSanitizedName(strings.Replace(se.Type, ".", "", -1))
	}
	if se.ImageSlug == "" {
		se.ImageSlug = fmt.Sprintf("image-%s-%s", se.EffectType, se.Slug)
	}
	if se.Extra && !strings.HasSuffix(se.ImageSlug, "-2") {
		se.ImageSlug = fmt.Sprintf("%s-2", se.ImageSlug)
	} else if strings.HasSuffix(se.ImageSlug, "-2") {
		se.Extra = true
		if !strings.HasSuffix(se.Slug, "-2") {
			se.Slug = fmt.Sprintf("%s-2", se.Slug)
		}
	}
	se.WebsiteLink = fmt.Sprintf("/%ss/%s", se.EffectType, se.Slug)
	if strings.HasSuffix(se.WebsiteLink, "-2") {
		se.WebsiteLink = se.WebsiteLink[:len(se.WebsiteLink)-2]
	}
	return nil
}

var (
	statusEffectFromSkillAuraMore = regexp.MustCompile(`([^\]]*)\[([^\]]+)\]`)
	statusEffectFromSkillAura     = regexp.MustCompile(`\[([^\]]+)\]`)
	debuffs                       = map[string]bool{
		"HP Burn":               true,
		"Poison":                true,
		"Decrease DEF":          true,
		"Decrease ACC":          true,
		"Decrease SPD":          true,
		"Decrease ATK":          true,
		"Weaken":                true,
		"Sleep":                 true,
		"Provoke":               true,
		"Freeze":                true,
		"Block Cooldown Skills": true,
		"Bomb":                  true,
		"Stun":                  true,
		"Block Buffs":           true,
		"Revive on Death":       true,
		"Heal Reduction":        true,
		"Leech":                 true,
	}
	buffs = map[string]bool{
		"Increase C. RATE": true,
		"Shield":           true,
		"Ally Protection":  true,
		"Reflect Damage":   true,
		"Increase DEF":     true,
		"Increase SPD":     true,
		"Increase ATK":     true,
		"Continuous Heal":  true,
		"Counterattack":    true,
		"Unkillable":       true,
		"Block Debuffs":    true,
		"Block Damage":     true,
	}
	stats = map[string]bool{
		"ATK":          true,
		"DEF":          true,
		"HP":           true,
		"SPD":          true,
		"Enemy MAX HP": true,
	}
	buffDebuffRateExtraSlug = map[string]string{
		"Continuous Heal":  "15%",
		"Decrease DEF":     "60%",
		"Ally Protection":  "50%",
		"Increase ATK":     "50%",
		"Increase C. RATE": "30%",
		"Increase SPD":     "30%",
		"Increase DEF":     "60%",
		"Reflect Damage":   "30%",
		"Heal Reduction":   "100%",
		"Decrease ACC":     "50%",
		"Decrease ATK":     "50%",
		"Decrease SPD":     "30%",
		"Poison":           " 5%",
		"Weaken":           "25%",
	}
)

type StatusEffectList []*StatusEffect

func (sl StatusEffectList) Sort() {
	sort.SliceStable(sl, func(i, j int) bool {
		if sl[i].EffectType != sl[j].EffectType {
			return sl[i].EffectType < sl[j].EffectType
		}
		return sl[i].Slug < sl[j].Slug
	})
}

func (se StatusEffect) GetPageTitle() string { return se.Type }

func (se StatusEffect) GetPageSlug() string { return se.Slug }

func (_ StatusEffect) GetPageTemplate() string {
	return "page-templates/template-champion-generated.php"
}

func (se StatusEffect) GetParentPageID() int {
	switch se.EffectType {
	case "buff":
		return 5313
	case "debuff":
		return 5318
	default:
		panic("unknown")
	}
}

func (se StatusEffect) GetPageContent(input io.Reader, output io.Writer, extraData map[string]interface{}) error {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}
	rawTemplate, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	tmpl, err := template.New("status_effect").Funcs(funcMap).Parse(string(rawTemplate))
	if err != nil {
		return err
	}
	extraData["StatusEffect"] = se
	err = tmpl.Execute(output, extraData)
	return err
}

func (se StatusEffect) GetPageExcerpt() string { return se.RawDescription }

func (se *StatusEffect) GetPageExtraData(dataDirectory string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	champions, errChampions := fetchChampions(dataDirectory)
	if errChampions != nil {
		return nil, errChampions
	}

	statusList, errStatus := fetchStatusEffects(dataDirectory)
	if errStatus != nil {
		return nil, errStatus
	}
	potentialUpgrade := fmt.Sprintf("%s-2", se.Slug)
	if _, ok := statusList[potentialUpgrade]; ok {
		data["UpgradedVersionOfStatusEffect"] = statusList[potentialUpgrade]
	}

	mapChampions := map[string]*Champion{}
	championEffect := map[string]map[string]*StatusEffect{}
	for _, champion := range champions {
		for _, skill := range champion.Skills {
			for _, effect := range skill.Effects {
				if se.equals(effect) {
					mapChampions[champion.Slug] = champion
					if _, ok := championEffect[champion.Slug]; !ok {
						championEffect[champion.Slug] = map[string]*StatusEffect{}
					}
					championEffect[champion.Slug][effect.Slug] = effect
				}
			}
		}
	}
	matching := make([]*Champion, len(mapChampions))
	idx := 0
	for _, champion := range mapChampions {
		matching[idx] = champion
		idx++
	}
	championEffects := map[string][]*StatusEffect{}
	for championSlug, _championEffects := range championEffect {
		_internalEffects := make([]*StatusEffect, len(_championEffects))
		idx := 0
		for _, effect := range _championEffects {
			_internalEffects[idx] = effect
			idx++
		}
		championEffects[championSlug] = _internalEffects
	}
	data["AllEffects"] = statusList
	data["ChampionEffectsMap"] = championEffects

	sort.SliceStable(matching, func(i, j int) bool {
		if matching[i].Rarity == matching[j].Rarity {
			return matching[i].Name < matching[j].Name
		}
		return rarityToRank[matching[i].Rarity] > rarityToRank[matching[j].Rarity]
	})

	data["AvailableChampions"] = matching

	return data, nil
}

func (se *StatusEffect) equals(oth *StatusEffect) bool {
	return se.WebsiteLink == oth.WebsiteLink
}

func (se StatusEffect) LinkName() string {
	return se.Slug
}
