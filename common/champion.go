package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Champion struct {
	Name            string                    `json:"name"`
	Rarity          string                    `json:"rarity"`
	Element         string                    `json:"element"`
	Type            string                    `json:"type"`
	Rating          Rating                    `json:"rating"`
	Slug            string                    `json:"slug"`
	Characteristics map[int64]Characteristics `json:"characteristics"`
	Auras           []Aura                    `json:"auras"`
	Skills          []Skill                   `json:"skills"`
	Faction         Faction                   `json:"faction"`
	FactionSlug     string                    `json:"faction_slug"`
	WebsiteLink     string                    `json:"website_link"`
	ImageSlug       string                    `json:"image_slug"`
}

func (c *Champion) Sanitize() error {
	sanitizedName, errName := GetSanitizedName(c.Name)
	if errName != nil {
		return errName
	}
	c.Name = sanitizedName
	if c.Characteristics == nil {
		c.Characteristics = map[int64]Characteristics{
			60: {},
		}
	}
	if c.Auras == nil {
		c.Auras = make([]Aura, 0)
	}
	if c.Skills == nil {
		c.Skills = make([]Skill, 0)
	}
	c.Slug = c.LinkName()

	// faction
	faction := &c.Faction
	errFaction := faction.Sanitize()
	if errFaction != nil {
		return errFaction
	}
	c.Faction = *faction
	c.FactionSlug = faction.Slug

	// link
	c.WebsiteLink = fmt.Sprintf("/champions/%s/", c.Slug)

	c.ImageSlug = fmt.Sprintf("image-champion-%s", c.Slug)

	return nil
}

func (c Champion) Filename() string {
	return fmt.Sprintf("%s.json", c.LinkName())
}

func (c Champion) LinkName() string {
	return GetLinkNameFromSanitizedName(c.Name)
}

func GetSanitizedName(name string) (string, error) {
	newName, err := strconv.Unquote(fmt.Sprintf(`"%s"`, name))
	if err != nil {
		return "", err
	} else if newName != name {
		return "", fmt.Errorf("please change name of %s", name)
	}
	newName = regexpSlug.ReplaceAllString(newName, `$1`)
	newName = strings.Trim(newName, " ")
	return newName, nil
}

func GetLinkNameFromSanitizedName(name string) string {
	for _, part := range []string{
		`'`, `"`, ` `,
	} {
		name = strings.Replace(name, part, "-", -1)
	}
	name = strings.ToLower(name)
	return name
}

var (
	regexpSlug = regexp.MustCompile("^([a-zA-Z0-9' -]+).*")
)

func (c Champion) GetPageTitle() string { return c.Name }

func (c Champion) GetPageSlug() string { return c.Slug }

func (_ Champion) GetPageTemplate() string { return "page-templates/template-champions.php" }

func (_ Champion) GetParentPageID() int { return 29 }
