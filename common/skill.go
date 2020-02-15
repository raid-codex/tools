package common

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Skill struct {
	Passive        bool            `json:"passive"`
	Name           string          `json:"name"`
	RawDescription string          `json:"raw_description"`
	Slug           string          `json:"slug"`
	Effects        []*StatusEffect `json:"effects"`
	DamageBasedOn  []string        `json:"damaged_based_on"`
	GIID           string          `json:"giid"`
	Upgrades       []*SkillData    `json:"upgrades"`
	ImageSlug      string          `json:"image_slug"`
	SkillNumber    string          `json:"skill_number"`
}

func (s *Skill) Sanitize() error {
	s.Slug = GetLinkNameFromSanitizedName(s.Name)
	s.Effects = nil
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
	if s.GIID != "" {
		s.ImageSlug = fmt.Sprintf("%x", md5.Sum([]byte(s.GIID)))
	}
	return nil
}

var (
	decreaseTMRegexp = regexp.MustCompile(`(?i)(Decrea|deplete|depleting)([^.<])+Turn Meter`)
	increaseTMRegexp = regexp.MustCompile(`(?i)(Fill|boost|steal|resets)([^.<])+Turn Meter`)
	anyTMRegexp      = regexp.MustCompile(`(?i)Turn.*Meter`)
)

func (s *Skill) lookForTurnMeter() error {
	for _, effect := range s.Effects {
		if effect.Slug == "increase-turn-meter" || effect.Slug == "decrease-turn-meter" {
			return nil
		}
	}
	matched := false
	if decreaseTMRegexp.MatchString(s.RawDescription) || strings.Contains(s.RawDescription, "Turn Meters decreased") {
		matched = true
		s.Effects = append(s.Effects, &StatusEffect{
			EffectType: "battle_enhancement",
			Type:       "Decrease Turn Meter",
		})
	}
	if increaseTMRegexp.MatchString(s.RawDescription) || strings.Contains(s.RawDescription, "Turn Meter will be increased") {
		matched = true
		s.Effects = append(s.Effects, &StatusEffect{
			EffectType: "battle_enhancement",
			Type:       "Increase Turn Meter",
		})
	}
	if anyTMRegexp.MatchString(s.RawDescription) && !matched {
		for _, wl := range []string{"Turn Meter is", "a full Turn Meter", "or have their Turn Meter filled"} {
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

func getEffectsFromDescription(effects []*StatusEffect, damageBasedOn []string, rawDescription string) ([]*StatusEffect, []string, error) {
	currentEffects := map[string]*StatusEffect{}
	/*	for _, effect := range effects {
		currentEffects[effect.Type] = effect
	}*/
	basedOn := map[string]bool{}
	/*	for _, stat := range damageBasedOn {
		basedOn[stat] = true
	}*/
	extract := statusEffectFromSkillAuraMore.FindAllString(rawDescription, -1)
	for _, m := range extract {
		m2 := statusEffectFromSkillAura.FindAllStringSubmatch(m, -1)
		if len(m2) != 1 {
			panic(fmt.Sprintf("wtf %+v", m2))
		}
		val := translateEffect(m2[0][1])
		switch true {
		case strings.Contains(m, fmt.Sprintf("is under a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("is under [%s]", val)),
			strings.Contains(m, fmt.Sprintf("if the target has a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("under [%s]", val)):
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
			return nil, nil, fmt.Errorf("unknown thing: %s (%+v) -- %s", val, m2, rawDescription)
		}
		if v, ok := buffDebuffRateExtraSlug[val]; ok && strings.Contains(m, v) {
			currentEffects[val].Extra = true
		}
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
