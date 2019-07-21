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
	Chance         float64 `json:"chance"`
	Turns          int64   `json:"turns"`
	Target         *Target `json:"target"`
	ImageSlug      string  `json:"image_slug"`
	Slug           string  `json:"slug"`
	WebsiteLink    string  `json:"website_link"`
	Extra          bool    `json:"extra"`
	RawDescription string  `json:"raw_description"`
	PlacesIf       string  `json:"places_if"`
	Amount         int64   `json:"amount"`
}

func (se *StatusEffect) Sanitize() error {
	if se.Slug == "" {
		se.Slug = GetLinkNameFromSanitizedName(strings.Replace(se.Type, ".", "", -1))
	}
	if se.ImageSlug == "" || se.Type == "Revive on Death" {
		se.ImageSlug = fmt.Sprintf("image-%s-%s", GetLinkNameFromSanitizedName(se.EffectType), se.Slug)
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
	if se.Target != nil {
		errTarget := se.Target.Sanitize()
		if errTarget != nil {
			return errTarget
		}
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
		"Heal Reduction":        true,
		"Leech":                 true,
	}
	buffs = map[string]bool{
		"Increase C. RATE":   true,
		"Shield":             true,
		"Ally Protection":    true,
		"Reflect Damage":     true,
		"Increase DEF":       true,
		"Increase SPD":       true,
		"Increase ATK":       true,
		"Continuous Heal":    true,
		"Counterattack":      true,
		"Unkillable":         true,
		"Block Debuffs":      true,
		"Block Damage":       true,
		"Increase C. DAMAGE": true,
		"Revive on Death":    true,
	}
	battleEnhancements = map[string]bool{
		"Ignore Block Damage":                 true,
		"Ignore Shield":                       true,
		"Critical Strike":                     true,
		"Increase Turn Meter":                 true,
		"Decrease Turn Meter":                 true,
		"Remove ALL Debuffs":                  true,
		"Heal":                                true,
		"Extra Turn":                          true,
		"Heal per DMG":                        true,
		"Revive":                              true,
		"Reset ALL Cooldowns":                 true,
		"Increase DMG":                        true,
		"DMG Reduction":                       true,
		"Transfer DMG":                        true,
		"Steal 1 Buff":                        true,
		"Transfer 1 Debuff":                   true,
		"Increase DMG per Debuff":             true,
		"Swap HP":                             true,
		"Extra Crit Chance":                   true,
		"Immune STUN":                         true,
		"Immune Freeze":                       true,
		"Immune Sleep":                        true,
		"Always Crit":                         true,
		"Extra Hit":                           true,
		"Remove 1 Buff":                       true,
		"Extra Crit DMG":                      true,
		"Crit Chance":                         true,
		"Ignore DEF":                          true,
		"ATK all Enemies":                     true,
		"Increase DMG per HP lost":            true,
		"Shield per Champ Max HP":             true,
		"DMG per Max HP":                      true,
		"ATK with Surplus DMG":                true,
		"Unlock Skypiercer":                   true,
		"Increase DMG per Buff":               true,
		"Ignore Block DMG":                    true,
		"Stack Damage upto X4":                true,
		"Stack Damage upto X5":                true,
		"Stack Damage upto X6":                true,
		"Stack Damage upto X7":                true,
		"Stack Damage upto X8":                true,
		"Stack Damage upto X9":                true,
		"Repeat Attack":                       true,
		"Remove 2 Buffs":                      true,
		"Shield per DMG":                      true,
		"Damage Per HP":                       true,
		"Remove 1 Debuff":                     true,
		"Deal ALL poison DMG":                 true,
		"Increase Debuffs":                    true,
		"Steal ALL Buffs":                     true,
		"Block Revive":                        true,
		"Decrease Bomb":                       true,
		"Increase Buff":                       true,
		"Decrease All Buffs":                  true,
		"Increase Debuff":                     true,
		"Heal per debuff":                     true,
		"DMG per HP":                          true,
		"Decrease DMG per Hit":                true,
		"Heal per DMG+Dead Ally":              true,
		"Increase DMG per Max HP + Dead Ally": true,
		"Increase All Buffs":                  true,
		"Remove ALL Buffs":                    true,
		"ATK With 2 Allies":                   true,
		"Decrease DMG":                        true,
		"DMG per Surplus Heal":                true,
		"Decrease MAX HP":                     true,
		"Remove CC Debuffs":                   true,
		"Increase DMG per Enemy Buff":         true,
		"Decrease Buff":                       true,
		"Deal all poison DMG":                 true,
		"Reset Skill Cooldown":                true,
		"Decrease Skill Cooldown":             true,
		"Deal Half of DMG":                    true,
		"Transfer ALL Debuffs":                true,
		"Increase ALL Debuffs":                true,
		"Crack Armor Skill":                   true,
		"Increase Meter per Stolen":           true,
		"Rec 5% DMG per Ally Alive":           true,
		"Spread 4 Debuffs":                    true,
		"Unlock Peril Skill":                  true,
		"Steal 2 Buffs":                       true,
		"Equalise all HP":                     true,
		"Immune Stun":                         true,
		"Increase RES per Buff":               true,
		"Self DMG":                            true,
		"Copy 1 Debuff":                       true,
		"Increase DMG as Enemy HP Dec":        true,
		"Increase DMG per Enemy Max HP":       true,
		"Increase DMG per Buff Removed":       true,
		"Reset Cooldown on Holy Sword":        true,
		"Extra DMG":                           true,
		"Increase DMG per Max HP":             true,
		"ATK":                                 true,
		"Reset Juggernaut Cooldown":           true,
		"Detonate Bombs":                      true,
		"Extra Crit Hit":                      true,
		"Decrease Cooldowns":                  true,
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

const (
	StatusEffect_CounterAttack = "counterattack"
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

func (se *StatusEffect) GetPageContent_Templates(tmpl *template.Template, output io.Writer, extraData map[string]interface{}) error {
	extraData["StatusEffect"] = se
	return tmpl.Execute(output, extraData)
}

func (se StatusEffect) GetPageExcerpt() string { return se.RawDescription }

func (se *StatusEffect) GetPageExtraData(dataDirectory string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	statusList, errStatus := fetchStatusEffects(dataDirectory)
	if errStatus != nil {
		return nil, errStatus
	}
	potentialUpgrade := fmt.Sprintf("%s-2", se.Slug)
	data["DescriptionClass"] = "col-xs-12"
	if _, ok := statusList[potentialUpgrade]; ok {
		data["UpgradedVersionOfStatusEffect"] = statusList[potentialUpgrade]
		data["DescriptionClass"] = "col-xs-12 col-md-6"
	}

	mapChampions := map[string]*Champion{}
	championEffect := map[string]map[string]*StatusEffect{}
	cl1, errCL := GetChampions(FilterChampionStatusEffect(se.Slug))
	if errCL != nil {
		return nil, errCL
	}
	cl2 := make(ChampionList, 0)
	if data["UpgradedVersionOfStatusEffect"] != nil {
		cl2, errCL = GetChampions(FilterChampionStatusEffect(data["UpgradedVersionOfStatusEffect"].(*StatusEffect).Slug))
		if errCL != nil {
			return nil, errCL
		}
	}
	runner := func(cl ChampionList, statusEffect *StatusEffect) {
		for _, champion := range cl {
			mapChampions[champion.Slug] = champion
			if _, ok := championEffect[champion.Slug]; !ok {
				championEffect[champion.Slug] = map[string]*StatusEffect{}
			}
			championEffect[champion.Slug][statusEffect.Slug] = statusEffect
		}
	}
	runner(cl1, se)
	if data["UpgradedVersionOfStatusEffect"] != nil {
		runner(cl2, data["UpgradedVersionOfStatusEffect"].(*StatusEffect))
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
