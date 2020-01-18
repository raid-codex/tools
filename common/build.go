package common

import (
	"fmt"
	"reflect"
	"sort"
)

type Build struct {
	From      string         `json:"from"`
	Author    string         `json:"author"`
	Locations []string       `json:"locations"`
	Sets      []string       `json:"sets"`
	Stats     *StatsPriority `json:"stats"`
}

func (b *Build) Set(piece string, stat *StatPriority) {
	if b.Stats == nil {
		b.Stats = &StatsPriority{}
	}
	v := reflect.Indirect(reflect.ValueOf(b.Stats))
	v.FieldByName(piece).Set(reflect.ValueOf(stat))
}

type StatsPriority struct {
	Weapon     *StatPriority `json:"weapon"`
	Helmet     *StatPriority `json:"helmet"`
	Shield     *StatPriority `json:"shield"`
	Gauntlets  *StatPriority `json:"gauntlets"`
	Chestplate *StatPriority `json:"chestplate"`
	Boots      *StatPriority `json:"boots"`
	Ring       *StatPriority `json:"ring"`
	Amulet     *StatPriority `json:"amulet"`
	Banner     *StatPriority `json:"banner"`
}

type StatPriority struct {
	MainStat        string   `json:"main_stat"`
	AdditionalStats []string `json:"additional_stats"`
}

var possibleStatsPerPiece = map[string]struct {
	Main       []string
	Additional []string
}{
	"Weapon": {
		Main:       []string{"ATK"},
		Additional: []string{"HP", "HP%", "ATK%", "SPD", "C.RATE", "C.DMG", "RESIST", "ACC"},
	},
	"Helmet": {
		Main:       []string{"HP"},
		Additional: []string{"HP%", "ATK", "ATK%", "SPD", "C.RATE", "C.DMG", "RESIST", "ACC"},
	},
	"Shield": {
		Main:       []string{"DEF"},
		Additional: []string{"HP", "HP%", "ATK", "ATK%", "SPD", "C.RATE", "C.DMG", "RESIST", "ACC"},
	},
	"Gauntlets": {
		Main: []string{"C.RATE", "C.DMG", "HP", "HP%", "ATK", "ATK%", "DEF", "DEF%"},
	},
}

func (ssp *StatsPriority) Sanitize() error {
	v := reflect.Indirect(reflect.ValueOf(ssp))
	for i := 0; i < v.NumField(); i++ {
		if err := v.Field(i).Interface().(*StatPriority).Sanitize(); err != nil {
			return err
		}
	}
	return nil
}

func (sp *StatPriority) Sanitize() error {
	n := make([]string, 0)
	for _, stat := range sp.AdditionalStats {
		if stat == sp.MainStat {
			continue
		}
		n = append(n, stat)
	}
	sp.AdditionalStats = n
	return nil
}

func (b *Build) Sanitize() error {
	sort.SliceStable(b.Locations, func(i, j int) bool { return b.Locations[i] < b.Locations[j] })
	sort.SliceStable(b.Sets, func(i, j int) bool { return b.Sets[i] < b.Sets[j] })
	if b.Stats != nil {
		return b.Stats.Sanitize()
	}
	if b.Sets == nil || len(b.Sets) == 0 {
		return fmt.Errorf("empty set for build")
	}
	return nil
}

func (b *Build) IsSameThan(oth *Build) bool {
	if oth.From == b.From && oth.Author == b.Author && len(oth.Locations) == len(b.Locations) {
		for idx := range oth.Locations {
			if b.Locations[idx] != oth.Locations[idx] {
				return false
			}
		}
		return true
	}
	return false
}
