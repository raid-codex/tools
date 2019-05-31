package common

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type StatusEffect struct {
	EffectType  string  `json:"effect_type"`
	Type        string  `json:"type"`
	Value       float64 `json:"value"`
	ImageSlug   string  `json:"image_slug"`
	Slug        string  `json:"slug"`
	WebsiteLink string  `json:"website_link"`
	Extra       bool    `json:"extra"`
}

func (se *StatusEffect) Sanitize() error {
	if se.Slug == "" {
		se.Slug = GetLinkNameFromSanitizedName(strings.Replace(se.Type, ".", "", -1))
	}
	se.ImageSlug = fmt.Sprintf("image-%s-%s", se.EffectType, se.Slug)
	if se.Extra && !strings.HasSuffix(se.ImageSlug, "-2") {
		se.ImageSlug = fmt.Sprintf("%s-2", se.ImageSlug)
	} else if strings.HasSuffix(se.ImageSlug, "-2") {
		se.Extra = true
	}
	se.WebsiteLink = fmt.Sprintf("/%s/%s", se.EffectType, se.Slug)
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
		"Debuff Spread":         true,
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
		"Block Heal":       "100%",
		"Decrease ACC":     "50%",
		"Decrease ATK":     "50%",
		"Decrease SPD":     "30%",
		"Poison":           "5%",
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
