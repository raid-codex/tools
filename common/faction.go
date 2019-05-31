package common

import (
	"fmt"
	"io"

	"github.com/raid-codex/tools/seo"
)

type Faction struct {
	Name               string   `json:"name"`
	Slug               string   `json:"slug"`
	WebsiteLink        string   `json:"website_link"`
	ImageSlug          string   `json:"image_slug"`
	NumberOfChampions  int64    `json:"number_of_champions"`
	DefaultDescription string   `json:"default_description"`
	SEO                *seo.SEO `json:"seo"`
}

func (f *Faction) Sanitize() error {
	name, err := GetSanitizedName(f.Name)
	if err != nil {
		return err
	}
	f.Name = name
	f.Slug = GetLinkNameFromSanitizedName(f.Name)
	f.WebsiteLink = fmt.Sprintf("/factions/%s/", f.Slug)
	f.ImageSlug = fmt.Sprintf("image-faction-%s", f.Slug)
	if f.NumberOfChampions > 0 {
		f.DefaultDescription = fmt.Sprintf(
			"%s is a faction from RAID Shadow Legends composed of %d champions",
			f.Name,
			f.NumberOfChampions,
		)
	} else {
		f.DefaultDescription = fmt.Sprintf(
			"%s is a faction from RAID Shadow Legends",
			f.Name,
		)
	}
	if f.SEO == nil {
		f.DefaultSEO()
	}

	return nil
}

func (f *Faction) DefaultSEO() {
	f.SEO = &seo.SEO{
		Title:       "%%title%% %%page%% %%sep%% %%parent_title%% %%sep%% %%sitename%%",
		Description: fmt.Sprintf("%s. Find out more on this Raid Shadow Legends codex.", f.DefaultDescription),
		Keywords: []string{
			"raid", "shadow", "legends", "factions", f.Name, f.Slug,
		},
	}
}

func (f Faction) Filename() string {
	return fmt.Sprintf("%s.json", f.LinkName())
}

func (f Faction) LinkName() string {
	return f.Slug
}

func (f Faction) GetPageTitle() string { return f.Name }

func (f Faction) GetPageSlug() string { return f.Slug }

func (_ Faction) GetPageTemplate() string { return "page-templates/template-faction.php" }

func (_ Faction) GetParentPageID() int { return 1730 }

func (_ Faction) GetPageContent(input io.Reader, output io.Writer, extraData map[string]interface{}) error {
	return nil
}

func (f Faction) GetPageExcerpt() string { return f.DefaultDescription }

func (c *Faction) GetPageExtraData(dataDirectory string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
