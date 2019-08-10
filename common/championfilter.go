package common

type ChampionFilter func(*Champion) bool

func FilterChampionNotSlug(slug string) ChampionFilter {
	return func(champion *Champion) bool {
		return champion.Slug != slug
	}
}

func FilterChampionFactionSlug(slug string) ChampionFilter {
	return func(champion *Champion) bool {
		return champion.FactionSlug == slug
	}
}

func FilterChampionStatusEffect(effectSlug string) ChampionFilter {
	return func(champion *Champion) bool {
		for _, skill := range champion.Skills {
			for _, effect := range skill.Effects {
				if effect.Slug == effectSlug {
					return true
				}
			}
			for _, upgrade := range skill.Upgrades {
				for _, effect := range upgrade.Effects {
					if effect.Slug == effectSlug {
						return true
					}
				}
			}
		}
		return false
	}
}

func FilterChampionStatusEffectOnSkill(skillNumber, effectSlug string) ChampionFilter {
	return func(champion *Champion) bool {
		for _, skill := range champion.Skills {
			if skill.SkillNumber != skillNumber {
				continue
			}
			for _, effect := range skill.Effects {
				if effect.Slug == effectSlug {
					return true
				}
			}
			for _, upgrade := range skill.Upgrades {
				for _, effect := range upgrade.Effects {
					if effect.Slug == effectSlug {
						return true
					}
				}
			}
		}
		return false
	}
}

func FilterChampionStatusEffectWithTargets(effectSlug string, targets ...string) ChampionFilter {
	return func(champion *Champion) bool {
		for _, skill := range champion.Skills {
			for _, effect := range skill.Effects {
				if effect.Slug == effectSlug && effect.Target != nil {
					for _, target := range targets {
						if effect.Target.Who == target {
							return true
						}
					}
				}
			}
			for _, upgrade := range skill.Upgrades {
				for _, effect := range upgrade.Effects {
					if effect.Slug == effectSlug && effect.Target != nil {
						for _, target := range targets {
							if effect.Target.Who == target {
								return true
							}
						}
					}
				}
			}
		}
		return false
	}
}

func FilterChampionSlug(slug string) ChampionFilter {
	return func(champion *Champion) bool { return champion.Slug == slug }
}
