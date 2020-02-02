package common

import (
	"fmt"
	"reflect"
)

type Rating struct {
	Overall       string `json:"overall"`
	Campaign      string `json:"campaign"`
	ArenaOff      string `json:"arena_offense"`
	ArenaDef      string `json:"arena_defense"`
	ClanBossWoGS  string `json:"clan_boss_without_giant_slayer"`
	ClanBosswGS   string `json:"clan_boss_with_giant_slayer"`
	IceGuardian   string `json:"ice_guardian"`
	Dragon        string `json:"dragon"`
	Spider        string `json:"spider"`
	FireKnight    string `json:"fire_knight"`
	Minotaur      string `json:"minotaur"`
	ForceDungeon  string `json:"force_dungeon"`
	MagicDungeon  string `json:"magic_dungeon"`
	SpiritDungeon string `json:"spirit_dungeon"`
	VoidDungeon   string `json:"void_dungeon"`
	FactionWars   string `json:"faction_wars"`
}

var (
	allowedRatings = map[string]bool{
		"":   true,
		"A":  true,
		"B":  true,
		"C":  true,
		"D":  true,
		"S":  true,
		"SS": true,
	}
	rankToInt = map[string]int{
		"D":  0,
		"C":  1,
		"B":  2,
		"A":  3,
		"S":  4,
		"SS": 5,
	}
	intToRank = map[int]string{
		5: "SS",
		4: "S",
		3: "A",
		2: "B",
		1: "C",
		0: "D",
	}
)

func (r *Rating) Sanitize() error {
	v := reflect.ValueOf(r)
	indV := reflect.Indirect(v)
	for i := 0; i < indV.NumField(); i++ {
		field := indV.Field(i)
		value := field.String()
		if _, ok := allowedRatings[value]; !ok {
			return fmt.Errorf("unknown rating %s", value)
		}
	}
	if r.Overall == "" {
		r.computeOverall()
	}
	return nil
}

func (r *Rating) computeOverall() {
	v := reflect.ValueOf(r)
	indV := reflect.Indirect(v)
	total := 0
	divideBy := 0
	for i := 0; i < indV.NumField(); i++ {
		if tag := indV.Type().Field(i).Tag.Get("json"); tag == "overall" {
			continue
		}
		value := indV.Field(i).String()
		if _, ok := rankToInt[value]; !ok {
			continue
		}
		total += rankToInt[value]
		divideBy += 1
	}
	if divideBy > 0 {
		r.Overall = overallRatioToRank(float32(total) / float32(divideBy))
	}
}

func overallRatioToRank(val float32) string {
	switch true {
	case val <= 0.5:
		return "D"
	case val <= 1.5:
		return "C"
	case val <= 2.5:
		return "B"
	case val <= 3.5:
		return "A"
	case val <= 4.5:
		return "S"
	case 4.5 < val:
		return "SS"
	default:
		return ""
	}
}

type Review struct {
	NumberOfReviews int64   `json:"amount"`
	Campaign        float64 `json:"campaign"`
	ArenaOff        float64 `json:"arena_offense"`
	ArenaDef        float64 `json:"arena_defense"`
	ClanBoss        float64 `json:"clan_boss"`
	IceGuardian     float64 `json:"ice_guardian"`
	Dragon          float64 `json:"dragon"`
	Spider          float64 `json:"spider"`
	FireKnight      float64 `json:"fire_knight"`
	Minotaur        float64 `json:"minotaur"`
	ForceDungeon    float64 `json:"force_dungeon"`
	MagicDungeon    float64 `json:"magic_dungeon"`
	SpiritDungeon   float64 `json:"spirit_dungeon"`
	VoidDungeon     float64 `json:"void_dungeon"`
}
