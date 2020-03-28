package common

import (
	"fmt"
	"html/template"
	"io"
	"sort"
	"strings"
	"time"
)

type Fusion struct {
	DateAdded        string              `json:"date_added"`
	TimeStart        *time.Time          `json:"time_start"`
	TimeEnd          *time.Time          `json:"time_end"`
	Active           bool                `json:"active"`
	Limit            *int64              `json:"limit"`
	ChampionSlug     string              `json:"champion_slug"`
	Name             string              `json:"name"`
	Slug             string              `json:"slug"`
	Ingredients      []*FusionIngredient `json:"ingredients"`
	ParentFusionSlug *string             `json:"parent_fusion_slug"`
	Schedule         *FusionSchedule     `json:"schedule"`
}

type FusionList []*Fusion

func (fl FusionList) Sort() {
	sort.SliceStable(fl, func(i, j int) bool {
		return fl[i].Slug < fl[j].Slug
	})
}

type FusionIngredient struct {
	ChampionSlug  string  `json:"champion_slug"`
	Level         int64   `json:"level"`
	Stars         int64   `json:"stars"`
	AscendedStars int64   `json:"ascended_stars"`
	Fusion        *Fusion `json:"fusion"`
	FusionSlug    *string `json:"fusion_slug"`
}

func (f *Fusion) Sanitize() error {
	if !strings.HasPrefix(f.Slug, "fusion-") {
		f.Slug = fmt.Sprintf("fusion-%s", f.Slug)
	}
	if f.Slug != fmt.Sprintf("fusion-%s", f.ChampionSlug) {
		parentSlug := strings.Replace(f.Slug, fmt.Sprintf("-%s", f.ChampionSlug), "", -1)
		f.ParentFusionSlug = &parentSlug
	}
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
		if ingredient.FusionSlug != nil {
			fusions, errFusion := GetFusions(FilterFusionSlug(*ingredient.FusionSlug))
			if errFusion != nil {
				return errFusion
			} else if len(fusions) != 1 {
				return fmt.Errorf("found %d fusions for %s", len(fusions), *ingredient.FusionSlug)
			} else {
				ingredient.Fusion = fusions[0]
			}
		}
	}
	if f.TimeStart != nil {
		f.Active = time.Now().After(*f.TimeStart)
		if f.TimeEnd != nil && f.Active {
			f.Active = time.Now().Before(*f.TimeEnd)
		}
	} else if f.TimeStart == nil && f.TimeEnd == nil {
		// unlimited
		f.Active = true
	}
	if f.Schedule != nil {
		f.Schedule.DateStart = f.TimeStart.Format("2006-01-02")
		f.Schedule.DateEnd = f.TimeEnd.Format("2006-01-02")
		if err := f.Schedule.Sanitize(); err != nil {
			return err
		}
	}
	return nil
}

func (f Fusion) GetPageTitle() string { return fmt.Sprintf("Fusion - %s", f.Name) }

func (f Fusion) GetPageSlug() string { return f.Slug }

func (_ Fusion) GetPageTemplate() string { return "page-templates/template-champion-generated.php" }

func (_ Fusion) GetParentPageID() int { return 11318 }

func (f Fusion) GetPageExcerpt() string { return f.Name }
func (f Fusion) LinkName() string       { return f.Slug }

func (f Fusion) GetPageContent(r io.Reader, output io.Writer, extraData map[string]interface{}) error {
	return fmt.Errorf("not implemented")
}

func (f *Fusion) GetPageContent_Templates(tmpl *template.Template, output io.Writer, extraData map[string]interface{}) error {
	extraData["Fusion"] = f
	return tmpl.Execute(output, extraData)
}

func (f *Fusion) hasChampion(championSlug string) bool {
	if f.ChampionSlug == championSlug {
		return true
	}
	for _, ingredient := range f.Ingredients {
		if ingredient.ChampionSlug == championSlug {
			return true
		}
		if ingredient.Fusion != nil && ingredient.Fusion.hasChampion(championSlug) {
			return true
		}
	}
	return false
}

func (f *Fusion) hasChampionAsIngredient(championSlug string) bool {
	for _, ingredient := range f.Ingredients {
		if ingredient.ChampionSlug == championSlug {
			return true
		}
	}
	return false
}

func (f *Fusion) GetPageExtraData(dataDirectory string) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	data["FusionLevel"] = 1
	champions, errChampions := GetChampions(func(c *Champion) bool {
		return f.hasChampion(c.Slug)
	})
	if errChampions != nil {
		return nil, errChampions
	}
	m := map[string]*Champion{}
	for _, champion := range champions {
		m[champion.Slug] = champion
	}
	data["Champions"] = m
	return data, nil
}
