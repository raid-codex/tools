package parse_static_data

type Element uint8

const (
	Element_Spirit Element = 3
)

type Rarity uint8

const (
	Rarity_Epic Rarity = 4
)

type Role uint8

const (
	Role_Def Role = 1
)

type AwakenMaterial string

const (
	AwakenMaterial_LesserMagic    AwakenMaterial = "103"
	AwakenMaterial_GreaterMagic   AwakenMaterial = "103"
	AwakenMaterial_SuperiorMagic  AwakenMaterial = "103"
	AwakenMaterial_LesserForce    AwakenMaterial = "111"
	AwakenMaterial_GreaterForce   AwakenMaterial = "112"
	AwakenMaterial_SuperiorForce  AwakenMaterial = "113"
	AwakenMaterial_LesserSpirit   AwakenMaterial = "121"
	AwakenMaterial_GreaterSpirit  AwakenMaterial = "122"
	AwakenMaterial_SuperiorSpirit AwakenMaterial = "123"
	AwakenMaterial_LesserVoid     AwakenMaterial = "131"
	AwakenMaterial_GreaterVoid    AwakenMaterial = "132"
	AwakenMaterial_SuperiorVoid   AwakenMaterial = "133"
	AwakenMaterial_LesserArcane   AwakenMaterial = "141"
	AwakenMaterial_GreaterArcane  AwakenMaterial = "142"
	AwakenMaterial_SuperiorArcane AwakenMaterial = "143"
)

type StatKind uint8

const (
	StatKind_HP  StatKind = 1
	StatKind_ATK StatKind = 2
)

type HeroType struct {
	ID   int64 `json:"Id"`
	Name struct {
		Key          string
		DefaultValue string
	}
	AvatarName      string
	ModelName       string
	Element         Element
	Role            Role
	Rarity          Rarity
	AwakenMaterials struct {
		RawValues map[AwakenMaterial]uint8
	}
	BaseStats struct {
		Health         int64
		Attack         int64
		Defense        int64 `json:"Defence"`
		Speed          int64
		Resistance     int64
		Accuracy       int64
		CriticalChance int64
		CriticalDamage int64
		CriticalHeal   int64
	}
	// SkillTypeIDs holds champion's skill IDs. Skill can evolve depending on the awakening level of the champion
	SkillTypeIDs        []int64 `json:"SkillTypeIds"`
	SummonWeight        int64
	IsLocationOnly      uint8
	Brain               int64
	ArtifactSuggestions []struct {
		KindID         int64 `json:"KindId"`
		SetKindID      int64 `json:"SetKindId"`
		PrimaryBonuses []struct {
			KindID     int64 `json:"KindId"`
			IsAbsolute uint8
		}
	}
	LeaderSkill struct {
		StatKindID StatKind `json:"StatKindId"`
		IsAbsolute uint8
		Amount     int64
	}
}
