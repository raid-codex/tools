package common

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/raid-codex/tools/seo"
)

type Champion struct {
	Name               string                    `json:"name"`
	Rarity             string                    `json:"rarity"`
	Element            string                    `json:"element"`
	Type               string                    `json:"type"`
	Rating             Rating                    `json:"rating"`
	Slug               string                    `json:"slug"`
	Characteristics    map[int64]Characteristics `json:"characteristics"`
	Auras              []Aura                    `json:"auras"`
	Skills             []Skill                   `json:"skills"`
	Faction            Faction                   `json:"faction"`
	FactionSlug        string                    `json:"faction_slug"`
	WebsiteLink        string                    `json:"website_link"`
	ImageSlug          string                    `json:"image_slug"`
	SEO                *seo.SEO                  `json:"seo"`
	DefaultDescription string                    `json:"default_description"`
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

	if c.Element != "" {
		c.DefaultDescription = fmt.Sprintf(
			"%s is a %s %s champion from the faction %s doing %s damage",
			c.Name,
			strings.ToLower(c.Rarity),
			strings.ToLower(c.Type),
			c.Faction.Name,
			strings.ToLower(c.Element),
		)
	} else {
		c.DefaultDescription = fmt.Sprintf(
			"%s is a %s %s champion from the faction %s",
			c.Name,
			strings.ToLower(c.Rarity),
			strings.ToLower(c.Type),
			c.Faction.Name,
		)
	}

	if c.SEO == nil {
		c.DefaultSEO()
	}

	c.defaultRating()

	return nil
}

func (c *Champion) defaultRating() {
	switch c.Rarity {
	case "Common", "Uncommon":
		// set ranking to D overall if no ranking
		if c.Rating.Overall == "" {
			c.Rating.Overall = "D"
		}
	}
}

func (c *Champion) DefaultSEO() {
	c.SEO = &seo.SEO{
		Title:       "%%title%% %%page%% %%sep%% %%parent_title%% %%sep%% %%sitename%%",
		Description: fmt.Sprintf("%s. Find out more on this Raid Shadow Legends codex.", c.DefaultDescription),
		Keywords: []string{
			"raid", "shadow", "legends", "champions", "tier", "list", c.Name, c.Slug,
		},
	}
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

type ChampionList []*Champion

func (cl ChampionList) Sort() {
	sort.SliceStable(cl, func(i, j int) bool {
		if cl[i].FactionSlug == cl[j].FactionSlug {
			if cl[i].Rarity == cl[j].Rarity {
				if cl[i].Rating.Overall == cl[j].Rating.Overall {
					return cl[i].Slug < cl[j].Slug
				}
				return ratingToRank[cl[i].Rating.Overall] > ratingToRank[cl[j].Rating.Overall]
			}
			return rarityToRank[cl[i].Rarity] > rarityToRank[cl[j].Rarity]
		}
		return cl[i].FactionSlug < cl[j].FactionSlug
	})
}

var (
	ratingToRank = map[string]int{
		"SS": 5,
		"S":  4,
		"A":  3,
		"B":  2,
		"C":  1,
		"D":  0,
	}
	rarityToRank = map[string]int{
		"Legendary": 4,
		"Epic":      3,
		"Rare":      2,
		"Uncommon":  1,
		"Common":    0,
	}
)
