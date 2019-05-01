package common

import (
	"fmt"
)

type Faction struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	WebsiteLink string `json:"website_link"`
	ImageSlug   string `json:"image_slug"`
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
	return nil
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
