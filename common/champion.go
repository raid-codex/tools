package common

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
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
	RecommendedBuilds  []Build                   `json:"recommended_builds"`
	Lore               string                    `json:"lore"`
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
	if c.RecommendedBuilds == nil {
		c.RecommendedBuilds = make([]Build, 0)
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

func (_ Champion) GetPageTemplate() string { return "page-templates/template-champion-generated.php" }

func (_ Champion) GetParentPageID() int { return 29 }

func (c Champion) GetPageExcerpt() string { return c.DefaultDescription }

func (c *Champion) GetPageContent(input io.Reader, output io.Writer) error {
	funcMap := template.FuncMap{
		"ToLower":      strings.ToLower,
		"DisplayGrade": grade,
		"Percentage":   func(s float64) int64 { return int64(s * 100.0) },
		"TrustAsHtml":  func(s string) template.HTML { return template.HTML(s) },
	}
	rawTemplate, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	tmpl, err := template.New("champion").Funcs(funcMap).Parse(string(rawTemplate))
	if err != nil {
		return err
	}

	err = tmpl.Execute(output, map[string]interface{}{
		"Champion": c,
	})
	return err
}

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

func grade(grade string) template.HTML {
	if grade == "" {
		return `<span class="champion-rating-none">No ranking yet</span>`
	}
	str := fmt.Sprintf(`<span class="champion-rating champion-rating-%s" title="%s">`, grade, gradeTitle(grade))
	for i := 0; i < 5; i++ {
		if i < ratingToRank[grade] {
			str += `<i class="fas fa-star"></i>`
		} else {
			str += `<i class="far fa-star"></i>`
		}
	}
	return template.HTML(str + `</span>`)
}

func gradeTitle(grade string) string {
	switch grade {
	case "D":
		return "not usable"
	case "C":
		return "viable"
	case "B":
		return "good"
	case "A":
		return "exceptional"
	case "S":
		return "top tier"
	case "SS":
		return "god tier"
	}
	return ""
}
