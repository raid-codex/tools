package parse_static_data

type SkillBonusType uint8

const (
	SkillBonusType_Damage             SkillBonusType = 0
	SkillBonusType_Buff_Debuff_Chance SkillBonusType = 2
	SkillBonusType_Cooldown           SkillBonusType = 3
)

type SkillType struct {
	ID       int64 `json:"Id"`
	Revision int64
	Name     struct {
		Key          string
		DefaultValue string
	}
	Description struct {
		Key          string
		DefaultValue string
	}
	Group                    int64
	Cooldown                 int64
	ReduceCooldownProhibited uint8
	IsHidden                 uint8
	ShowDamageScale          uint8
	Visibility               uint8
	SkillLevelBonuses        []struct {
		SkillBonusType SkillBonusType
		Value          WeirdValue
	}
	Effects []struct {
		ID           int64 `json:"Id"`
		KindID       int64 `json:"KindId"`
		Group        int64
		TargetParams struct {
			TargetType         uint8
			Exclusive          uint8
			FirstHitInSelected uint8
		}
		IsEffectDescription      uint8
		ConsidersDead            uint8
		LeaveThroughDeath        uint8
		DoesntSetSkillOnCooldown uint8
		IgnoresCooldown          uint8
		IsUnique                 uint8
		IterationChanceRolling   uint8
		Relation                 struct {
			EffectKindIDs         []int64 `json:"EffectKindIds"`
			Phase                 int64
			ActivateOnGlancingHit int64
		}
		Condition             string
		Count                 int64
		StackCount            int64
		MultiplierFormula     string
		ValueCap              string
		PersistsThroughRounds uint8
		SnapshotRequired      uint8
		HealParams            struct {
			CanBeCritical uint8
		}
	}
}
