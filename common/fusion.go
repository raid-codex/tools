package common

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Fusion struct {
	TimeStart       *time.Time          `json:"time_start"`
	TimeEnd         *time.Time          `json:"time_end"`
	Limit           *int64              `json:"limit"`
	ChampionSlug    string              `json:"champion_slug"`
	Name            string              `json:"name"`
	Slug            string              `json:"slug"`
	Ingredients     []*FusionIngredient `json:"ingredients"`
	ChildrenFusions []*Fusion           `json:"children_fusions"`
}

type FusionList []*Fusion

func (fl FusionList) Sort() {
	sort.SliceStable(fl, func(i, j int) bool {
		return fl[i].Slug < fl[j].Slug
	})
}

type FusionIngredient struct {
	ChampionSlug  string `json:"champion_slug"`
	Level         int64  `json:"level"`
	Stars         int64  `json:"stars"`
	AscendedStars int64  `json:"ascended_stars"`
}

func (f *Fusion) Sanitize() error {
	if !strings.HasPrefix(f.Slug, "fusion-") {
		f.Slug = fmt.Sprintf("fusion-%s", f.Slug)
	}
	newChildren := make([]*Fusion, len(f.ChildrenFusions))
	for idx, childFusion := range f.ChildrenFusions {
		if !strings.HasPrefix(childFusion.Slug, "fusion-") {
			childFusion.Slug = fmt.Sprintf("fusion-%s", childFusion.Slug)
		}
		// ensure children fusions exists
		fusions, errFusion := GetFusions(FilterFusionSlug(childFusion.Slug))
		if errFusion != nil {
			return errFusion
		} else if len(fusions) != 1 {
			return fmt.Errorf("found %d fusions with slug %s", len(fusions), childFusion.Slug)
		} else {
			newChildren[idx] = fusions[0]
		}
	}
	f.ChildrenFusions = newChildren
	// ensure champions mentioned exists
	{
		champions, errChampions := GetChampions(FilterChampionSlug(f.ChampionSlug))
		if errChampions != nil {
			return errChampions
		} else if len(champions) != 1 {
			return fmt.Errorf("found %d champions with slug %s", len(champions), f.ChampionSlug)
		}
	}
	for _, ingredient := range f.Ingredients {
		champions, errChampions := GetChampions(FilterChampionSlug(ingredient.ChampionSlug))
		if errChampions != nil {
			return errChampions
		} else if len(champions) != 1 {
			return fmt.Errorf("found %d champions with slug %s", len(champions), ingredient.ChampionSlug)
		}
	}
	return nil
}
