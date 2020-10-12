package common

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
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
	MainStat        string   `json:"main_stat,omitempty"`
	MainStats       []string `json:"main_stats"`
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

var (
	statPartsToRemove = []string{
		"Stat/Substat: ",
	}
	statReplacement = map[string]string{
		"Attack%":          "ATK%",
		"Critical Damage":  "C.DMG",
		"Critical Rate":    "C.RATE",
		"Speed":            "SPD",
		"Accuracy":         "ACC",
		"Resistance":       "RES",
		"Defense%":         "DEF%",
		"Defense":          "DEF",
		"Critical Rate%":   "C.RATE",
		"Attack":           "ATK",
		"Atack":            "ATK",
		"Resist":           "RES",
		"Crittical Damage": "C.DMG",
		"RESIST":           "RES",
		"Critical Damge":   "C.DMG",
		"Defense%%":        "DEF%",
		"C.RATE%":          "C.RATE",
		"Crit Rate":        "C.RATE",
	}
	knownStatErrs = map[string][]string{
		"HP% Critical Rate": []string{"HP%", "C.RATE"},
		"Accuracy/Resist":   []string{"ACC", "RES"},
		"HP% Speed":         []string{"HP%", "SPD"},
		"ATK% C.RATE":       []string{"ATK%", "C.RATE"},
		"C.DMG SPD":         []string{"C.DMG", "SPD"},
	}
)

func (sp *StatPriority) Sanitize() error {
	if sp.MainStat != "" {
		sp.MainStats = []string{}
		for _, str := range strings.Split(sp.MainStat, "/") {
			str = strings.Trim(str, " ")
			sp.MainStats = append(sp.MainStats, str)
		}
		sp.MainStat = "" // remove it
	}
	n := make([]string, 0)
	for _, stat := range sp.AdditionalStats {
		for _, toRemove := range statPartsToRemove {
			if strings.Contains(stat, toRemove) {
				stat = strings.Replace(stat, toRemove, "", 1)
			}
		}
		if stats, ok := knownStatErrs[stat]; ok {
			n = append(n, stats...)
		} else {
			n = append(n, stat)
		}
	}
	sp.AdditionalStats = n
	n = make([]string, 0)
loop_stat:
	for _, stat := range sp.AdditionalStats {
		for _, mainStat := range sp.MainStats {
			if mainStat == stat {
				continue loop_stat
			}
		}
		if v, ok := statReplacement[stat]; ok {
			stat = v
		}
		n = append(n, stat)
	}
	sp.AdditionalStats = n
	for idx, stat := range sp.MainStats {
		if strings.Contains(stat, "/") && stat != "N/A" {
			stats := strings.Split(stat, "/")
			stat = stats[0]
			sp.MainStats = append(sp.MainStats, stats[1:]...)
		}
		if v, ok := statReplacement[stat]; ok {
			stat = v
		}
		sp.MainStats[idx] = stat
	}
	if len(sp.MainStats) == 2 && sp.MainStats[0] == "N" && sp.MainStats[1] == "A" {
		sp.MainStats = []string{"N/A"}
	}
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
