package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/raid-codex/tools/utils"
)

type Aura struct {
	RawDescription string          `json:"raw_description"`
	Effects        []*StatusEffect `json:"effects"`
	Stats          []string        `json:"stats"`
	Locations      []string        `json:"locations"`
	Value          int64           `json:"value"`
	Percentage     bool            `json:"percentage"`
}

var (
	regexpAura       = regexp.MustCompile(`^Increases.+(?P<Stat>ATK|DEF|ACC|SPD|RES|C\.RATE|HP) in (?P<Location>.+) by (?P<Value>\d+%?)`)
	auraReplacements = []struct {
		lookFor   string
		replaceBy string
	}{
		{"Accuracy", "ACC"},
		{"Hit Point", "HP"},
		{"Attack", "ATK"},
		{"Defense", "DEF"},
		{"Speed", "SPD"},
		{"Critical Rate", "C.RATE"},
		{"C. RATE", "C.RATE"},
		{"Resist", "RES"},
		{"RESIST", "RES"},
		{"Arena Battles", "Arena"},
	}
)

func (a *Aura) Sanitize() error {
	a.Stats = make([]string, 0)
	a.Locations = make([]string, 0)
	for _, repl := range auraReplacements {
		a.RawDescription = strings.Replace(a.RawDescription, repl.lookFor, repl.replaceBy, -1)
	}
	matches := regexpAura.FindStringSubmatch(a.RawDescription)
	if matches == nil || len(matches) != 4 {
		return fmt.Errorf("regexp did not match anything in raw description '%s'", a.RawDescription)
	}
	a.Stats = append(a.Stats, matches[1])
	a.Locations = append(a.Locations, matches[2])
	if strings.HasSuffix(matches[3], "%") {
		a.Percentage = true
		matches[3] = matches[3][0 : len(matches[3])-1]
	} else {
		a.Percentage = false
	}
	value, err := strconv.ParseInt(matches[3], 10, 64)
	if err != nil {
		return err
	}
	a.Value = value
	effects, _, err := getEffectsFromDescription(a.Effects, []string{}, a.RawDescription)
	if err != nil {
		return err
	}
	a.Effects = effects
	a.Stats = utils.UniqueSlice(a.Stats)
	a.Locations = utils.UniqueSlice(a.Locations, GetLinkNameFromSanitizedName, ConvertLocation)
	return nil
}
