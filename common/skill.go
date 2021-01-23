package common

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type Skill struct {
	Passive        bool            `json:"passive"`
	Name           string          `json:"name"`
	RawDescription string          `json:"raw_description"`
	Slug           string          `json:"slug"`
	Effects        []*StatusEffect `json:"effects"`
	DamageBasedOn  []string        `json:"damaged_based_on"`
	GIID           string          `json:"giid"`
	Cooldown       int64           `json:"cooldown"`
	Upgrades       []*SkillData    `json:"upgrades"`
	ImageSlug      string          `json:"image_slug"`
	SkillNumber    string          `json:"skill_number"`
}

func (s *Skill) Sanitize() error {
	s.Slug = GetLinkNameFromSanitizedName(s.Name)
	s.Effects = nil
	if err := s.parseRawSkill(); err != nil {
		return err
	}
	effects, basedOn, err := getEffectsFromDescription(s.Effects, s.DamageBasedOn, s.RawDescription)
	if err != nil {
		return err
	}
	s.Effects = effects
	if err := s.lookForTurnMeter(); err != nil {
		return err
	}
	s.DamageBasedOn = basedOn
	for _, effect := range s.Effects {
		errSanitize := effect.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
	}
	if s.Upgrades == nil {
		s.Upgrades = make([]*SkillData, 0)
	}
	for _, upgrade := range s.Upgrades {
		errUpgrade := upgrade.Sanitize()
		if errUpgrade != nil {
			return errUpgrade
		}
	}
	sort.SliceStable(s.Upgrades, func(i, j int) bool {
		return s.Upgrades[i].Level < s.Upgrades[j].Level
	})
	// no cooldown yet, fetch from upgrades
	if s.Cooldown == 0 && !s.Passive && s.SkillNumber != "A1" && len(s.Upgrades) > 0 {
		s.Cooldown = s.Upgrades[0].Cooldown
	}
	if s.GIID != "" {
		s.ImageSlug = fmt.Sprintf("%x", md5.Sum([]byte(s.GIID)))
	}
	return nil
}

var (
	decreaseTMRegexp       = regexp.MustCompile(`(?i)(steals|Decrea|deplete|depleting)([^.<])+Turn Meter`)
	increaseTMRegexp       = regexp.MustCompile(`(?i)(steals|Fill|boost|steal|resets)([^.<])+Turn Meter`)
	anyTMRegexp            = regexp.MustCompile(`(?i)Turn.*Meter`)
	descriptionReplacement = regexp.MustCompile(`(\d+)\.(\d+)`)
)

func (s *Skill) lookForTurnMeter() error {
	for _, effect := range s.Effects {
		if effect.Slug == "increase-turn-meter" || effect.Slug == "decrease-turn-meter" {
			return nil
		}
	}
	// match description by removing potential "Steals 7.5% of Turn meter" things which invalidates decrease/increase regexp
	description := descriptionReplacement.ReplaceAllString(s.RawDescription, "$1,$2")
	matched := false
	if decreaseTMRegexp.MatchString(description) || strings.Contains(s.RawDescription, "Turn Meters decreased") {
		matched = true
		s.Effects = append(s.Effects, &StatusEffect{
			EffectType: "battle_enhancement",
			Type:       "Decrease Turn Meter",
		})
	}
	if increaseTMRegexp.MatchString(description) || strings.Contains(s.RawDescription, "Turn Meter will be increased") {
		matched = true
		s.Effects = append(s.Effects, &StatusEffect{
			EffectType: "battle_enhancement",
			Type:       "Increase Turn Meter",
		})
	}
	if anyTMRegexp.MatchString(s.RawDescription) && !matched {
		for _, wl := range []string{"Turn Meter is", "a full Turn Meter", "or have their Turn Meter filled", "Revives Skullsworn with 50% HP and 50% Turn Meter", "revives all dead allies with 50% HP and 50% Turn Meter", "with the highest Turn Meter", "Immune to Turn Meter decreasing effects.", "Revives a dead ally with 50% HP and 50% Turn Meter", "has more than 75% Turn Meter.", "Revives 2 Random allies with 20% HP and 20% Turn Meter"} {
			if strings.Contains(s.RawDescription, wl) {
				matched = true
			}
		}
	}
	if anyTMRegexp.MatchString(s.RawDescription) && !matched {
		return fmt.Errorf("matched Turn Meter in skill %s but did not find any increase/decrease\n%s", s.Name, s.RawDescription)
	}
	return nil
}

func (s *Skill) SetSkillData(sd *SkillData) {
	ns := make([]*SkillData, 0)
	ns = append(ns, sd)
	for _, skillData := range s.Upgrades {
		if skillData.Level != sd.Level {
			ns = append(ns, skillData)
		}
	}
	s.Upgrades = ns
}

func translateEffect(str string) string {
	if _, ok := effectTranslate[str]; ok {
		return effectTranslate[str]
	} else if strings.HasPrefix(str, "Dec ") {
		return strings.Replace(str, "Dec ", "Decrease ", -1)
	} else if strings.HasPrefix(str, "Inc ") {
		return strings.Replace(str, "Inc ", "Increase ", -1)
	}
	return str
}

var (
	effectTranslate = map[string]string{
		"Heal Red":                     "Heal Reduction",
		"Crit":                         "Critical Strike",
		"Inc Meter":                    "Increase Turn Meter",
		"Dec Meter":                    "Decrease Turn Meter",
		"Reflect DMG":                  "Reflect Damage",
		"Inc Crit Rate":                "Increase C. RATE",
		"Skills on Cooldown":           "Block Cooldown Skills",
		"Skill on Cooldown":            "Block Cooldown Skills",
		"Inc Skill Cooldown":           "Block Cooldown Skills",
		"Inc 1 Skill Cooldown":         "Block Cooldown Skills",
		"Increase Cooldowns":           "Block Cooldown Skills",
		"Dec Buffs":                    "Decrease Buff",
		"Decrease Max HP":              "Decrease MAX HP",
		"Heal per HP":                  "Heal",
		"Cont Heal":                    "Continuous Heal",
		"Counter":                      "Counterattack",
		"Ally Prot":                    "Ally Protection",
		"Block DMG":                    "Block Damage",
		"Fill Meter":                   "Increase Turn Meter",
		"cont heal":                    "Continuous Heal",
		"Empty Meter":                  "Decrease Turn Meter",
		"block debuffs":                "Block Debuffs",
		"Remove All Debuffs":           "Remove ALL Debuffs",
		"Shield per dmg":               "Shield per DMG",
		"heal":                         "Heal",
		"Burn":                         "HP Burn",
		"ignore block dmg":             "Ignore Block DMG",
		"Revive Block":                 "Block Revive",
		"transfer 1 debuff":            "Transfer 1 Debuff",
		"Inc DMG as HP lost":           "Increase DMG per HP lost",
		"reflect dmg":                  "Reflect Damage",
		"shield":                       "Shield",
		"Inc Crit DMG":                 "Increase C. DAMAGE",
		"skills on cooldown":           "Block Cooldown Skills",
		"Dec DEF per Self Debuff":      "Decrease DEF",
		"extra turn":                   "Extra Turn",
		"Inc Buffs":                    "Increase All Buffs",
		"Shield per DEF":               "Shield",
		"Shield per Max HP":            "Shield",
		"extra crit chance":            "Extra Crit Chance",
		"deal all poison dmg":          "Deal all poison DMG",
		"dec spd":                      "Decrease SPD",
		"Inc DMG per Debuff on Target": "Increase DMG per Debuff",
		"2x bomb":                      "Bomb",
		"ignore shield":                "Ignore Shield",
		"sleep":                        "Sleep",
		"steal 1 buff":                 "Steal 1 Buff",
		"Shield per LVL":               "Shield",
		"Heal per Self Max HP":         "Heal",
		"shield per max hp":            "Shield",
		"Reset all Cooldowns":          "Reset ALL Cooldowns",
		"Heal per ATK":                 "Heal",
		"Revive with Full HP":          "Revive",
		"Dec Max HP per DMG":           "Decrease MAX HP",
		"Dec MAX HP Per DMG":           "Decrease MAX HP",
		"Dec Max HP":                   "Decrease MAX HP",
		"Heal Per Surplus DMG":         "Heal",
		"weaken":                       "Weaken",
		"inc debuffs":                  "Increase Debuffs",
		"Block Debuff":                 "Block Debuffs",
		"Inc Meter per Hit":            "Increase Turn Meter",
		"Heal ALL to Highest HP":       "Heal",
		"shield per dmg":               "Shield",
		"Inc Meter per Debuff":         "Increase Turn Meter",
		"Heal per Debuff":              "Heal",
		"heal per debuff":              "Heal",
		"Shield per DEF+Crit DMG":      "Shield",
		"Dec Meter Per Buff":           "Decrease Turn Meter",
		"Poison x3":                    "Poison",
		"dec buff":                     "Decrease Buff",
		"Heal per Buff":                "Heal",
		"Inc DEF per Dead Ally":        "Increase DEF",
		"Deal ALL Poison":              "Deal all poison DMG",
		"Remove All Buffs":             "Remove ALL Buffs",
		"increase DEF":                 "Increase DEF",
		"Increase C. Rate":             "Increase C. RATE",
		"Increase C.RATE":              "Increase C. RATE",
		"Decrease C.DMG":               "Decrease C. DMG",
		"Increase C.DMG":               "Increase C. DMG",
	}
)

func (sd *SkillData) AddEffect(effect string, who string, turns int64, chance float64, placesIf string, value float64, amount int64) {
	if effect == "" {
		return
	}
	realEffect := translateEffect(effect)
	placesIf = translateEffect(placesIf)
	statusEffect := &StatusEffect{
		Type:     realEffect,
		Turns:    turns,
		Target:   &Target{Who: who},
		Chance:   chance,
		Value:    value,
		PlacesIf: placesIf,
		Amount:   amount,
	}
	if _, ok := debuffs[statusEffect.Type]; ok {
		statusEffect.EffectType = "debuff"
	} else if _, ok := buffs[statusEffect.Type]; ok {
		statusEffect.EffectType = "debuff"
	} else if _, ok := battleEnhancements[statusEffect.Type]; ok {
		statusEffect.EffectType = "battle_enhancement"
	} else {
		panic(fmt.Errorf("unknown effect '%s'", statusEffect.Type))
	}
	ne := make([]*StatusEffect, 0)
	ne = append(ne, statusEffect)
	for _, se := range sd.Effects {
		if se.Type != statusEffect.Type {
			ne = append(ne, se)
		}
	}
	sd.Effects = ne
}

func (s *Skill) getSentencesFromRawDescription() []string {
	split := strings.Split(s.RawDescription, "<br>")
	split2 := make([]string, 0)
	for _, s := range split {
		split3 := strings.Split(s, ".")
		for _, s3 := range split3 {
			if s3 == "" {
				continue
			}
			split2 = append(split2, s3)
		}
	}
	return split2
}

func (s *Skill) parseRawSkill() error {
	currentEffects := map[string]*StatusEffect{}
	for _, sentence := range s.getSentencesFromRawDescription() {
		if strings.HasPrefix(sentence, "Lvl") || strings.HasPrefix(sentence, "Level") {
			break
		}
		logrus.Debugf("sentence: %s\n", sentence)
		/*if err := parseStatusEffectsFromSkillDescription(sentence, currentEffects, basedOn); err != nil {
			return nil, err
		}*/
		if err := parseTargetsOfAttack(sentence); err != nil {
			return err
		}
	MAINLOOP:
		for _, matcher := range regexpMatchers {
			if matcher.Regexp.MatchString(sentence) {
				for _, negMatcher := range matcher.NegativeRegexps {
					if negMatcher.MatchString(sentence) {
						continue MAINLOOP
					}
				}
				currentEffects[matcher.Type] = &StatusEffect{
					EffectType: matcher.EffectType,
					Type:       matcher.Type,
				}
			}
		}
	}
	logrus.Debugf("%+v\n", currentEffects)
	for _, effect := range currentEffects {
		s.Effects = append(s.Effects, effect)
	}
	return nil
}

var (
	parseTargetsAttack = regexp.MustCompile(`Attacks (\d+|all) enem`)
	regexpMatchers     = []struct {
		Regexp          *regexp.Regexp
		EffectType      string
		Type            string
		NegativeRegexps []*regexp.Regexp
	}{
		{
			Regexp:     regexp.MustCompile(`([iI]ncreas|[eE]xtend).+the duration.+ debuff`),
			EffectType: "battle_enhancement",
			Type:       "Debuff extend",
		},
		{
			Regexp:     regexp.MustCompile(`([iI]ncreas|[eE]xtend).+the duration.+ buff`),
			EffectType: "battle_enhancement",
			Type:       "Buff extend",
			NegativeRegexps: []*regexp.Regexp{
				regexp.MustCompile(`([iI]ncreas|[eE]xtend).+the duration.+ debuff.+ under .+buff`),
			},
		},
		{
			Regexp:     regexp.MustCompile(` [hH]eal[^\]]`),
			EffectType: "battle_enhancement",
			Type:       "Heal",
		},
	}
)

func parseTargetsOfAttack(sentence string) error {
	matches := parseTargetsAttack.FindAllStringSubmatch(sentence, -1)
	for _, m := range matches {
		if len(m) != 2 {
			return fmt.Errorf("wtf %+v", m)
		}
		targets := m[1]
		logrus.Debugf("targets: %s\n", targets)
	}
	return nil
}

func parseStatusEffectsFromSkillDescription(sentence string, currentEffects map[string]*StatusEffect, basedOn map[string]bool) error {
	extract := statusEffectFromSkillAuraMore.FindAllString(sentence, -1)
	for _, m := range extract {
		m2 := statusEffectFromSkillAura.FindAllStringSubmatch(m, -1)
		if len(m2) != 1 {
			return fmt.Errorf("wtf %+v", m2)
		}
		val := translateEffect(m2[0][1])
		switch true {
		case strings.Contains(m, fmt.Sprintf("is under a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("is under [%s]", val)),
			strings.Contains(m, fmt.Sprintf("if the target has a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("under [%s]", val)),
			strings.Contains(m, fmt.Sprintf("under an [%s]", val)),
			strings.Contains(m, fmt.Sprintf("under a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("from all [%s]", val)),
			strings.Contains(m, fmt.Sprintf("removes any [%s]", val)),
			strings.Contains(m, fmt.Sprintf("an enemy places a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("except [%s]", val)),
			strings.Contains(m, fmt.Sprintf("an ally receives a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("Ignores [%s]", val)),
			strings.Contains(sentence, fmt.Sprintf("Ignores [Shield] and [%s]", val)), // special case for basileus roanas
			strings.Contains(m, fmt.Sprintf("Will ignore [%s]", val)),
			strings.Contains(sentence, fmt.Sprintf("Will ignore [Shield] and [%s]", val)): // special for crypt witch
			// this is pure garbage
			continue
		}
		if _, ok := debuffs[val]; ok {
			currentEffects[val] = &StatusEffect{
				EffectType: "debuff",
				Type:       val,
			}
		} else if _, ok := buffs[val]; ok {
			currentEffects[val] = &StatusEffect{
				EffectType: "buff",
				Type:       val,
			}
		} else if _, ok := stats[val]; ok {
			basedOn[val] = true
		} else if _, ok := battleEnhancements[val]; ok {
			currentEffects[val] = &StatusEffect{
				EffectType: "battle_enhancement",
				Type:       val,
			}
		} else {
			return fmt.Errorf("unknown thing: %s (%+v) -- %s", val, m2, sentence)
		}
		if v, ok := buffDebuffRateExtraSlug[val]; ok && strings.Contains(m, v) {
			currentEffects[val].Extra = true
		}
	}
	return nil
}

func getEffectsFromDescription(effects []*StatusEffect, damageBasedOn []string, rawDescription string) ([]*StatusEffect, []string, error) {
	currentEffects := map[string]*StatusEffect{}
	for _, effect := range effects {
		currentEffects[effect.Type] = effect
	}
	basedOn := map[string]bool{}
	if err := parseStatusEffectsFromSkillDescription(rawDescription, currentEffects, basedOn); err != nil {
		return nil, nil, err
	}
	newEffects := make([]*StatusEffect, len(currentEffects))
	idx := 0
	for _, effect := range currentEffects {
		errSanitize := effect.Sanitize()
		if errSanitize != nil {
			return nil, nil, errSanitize
		}
		newEffects[idx] = effect
		idx++
	}
	newBasedOn := make([]string, len(basedOn))
	idx = 0
	for stat, ok := range basedOn {
		if !ok {
			continue
		}
		newBasedOn[idx] = stat
		idx++
	}
	sort.SliceStable(newEffects, func(i, j int) bool {
		return newEffects[i].Slug < newEffects[j].Slug
	})
	sort.SliceStable(newBasedOn, func(i, j int) bool {
		return newBasedOn[i] < newBasedOn[j]
	})
	return newEffects, newBasedOn, nil
}
