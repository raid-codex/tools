package common

import (
	"encoding/json"
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
	Rating             *Rating                   `json:"rating"`
	Reviews            *Review                   `json:"reviews"`
	Slug               string                    `json:"slug"`
	Characteristics    map[int64]Characteristics `json:"characteristics"`
	Auras              []*Aura                   `json:"auras"`
	Skills             []*Skill                  `json:"skills"`
	Faction            Faction                   `json:"faction"`
	FactionSlug        string                    `json:"faction_slug"`
	WebsiteLink        string                    `json:"website_link"`
	ImageSlug          string                    `json:"image_slug"`
	SEO                *seo.SEO                  `json:"seo"`
	DefaultDescription string                    `json:"default_description"`
	RecommendedBuilds  []*Build                  `json:"recommended_builds"`
	Lore               string                    `json:"lore"`
	GIID               string                    `json:"giid"`
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
		c.Auras = make([]*Aura, 0)
	}
	if c.Skills == nil {
		c.Skills = make([]*Skill, 0)
	}
	if c.RecommendedBuilds == nil {
		c.RecommendedBuilds = make([]*Build, 0)
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

	for idx, skill := range c.Skills {
		if c.GIID != "" && skill.GIID == "" {
			skill.GIID = fmt.Sprintf("%s_s%d", c.GIID, idx+1)
		}
		errSanitize := skill.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
	}

	for _, aura := range c.Auras {
		errSanitize := aura.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
	}

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
		StructuredData: []json.RawMessage{},
	}
	championDoc := map[string]interface{}{
		"@context":    "http://schema.org/",
		"@type":       "Person",
		"name":        c.Name,
		"url":         fmt.Sprintf("https://raid-codex.com%s", c.WebsiteLink),
		"image":       fmt.Sprintf("https://raid-codex.com/wp-content/uploads/champions/%s.jpg", c.ImageSlug),
		"description": fmt.Sprintf("Member of the faction %s, %s is a champion of %s rarity and of %s type", c.Faction.Name, c.Name, c.Rarity, c.Type),
		"affiliation": map[string]interface{}{
			"@type":    "Organization",
			"@context": "http://schema.org",
			"name":     c.Faction.Name,
		},
	}
	rawMessage, _ := json.Marshal(championDoc)
	c.SEO.StructuredData = append(c.SEO.StructuredData, json.RawMessage(rawMessage))
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
		`'`, `"`, ` `, `_`, `(P)`, `+`, `%`,
	} {
		name = strings.Replace(name, part, "-", -1)
	}
	name = strings.ToLower(name)
	name = strings.Trim(name, " ")
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

func (c *Champion) GetPageContent(input io.Reader, output io.Writer, extraData map[string]interface{}) error {
	funcMap := template.FuncMap{
		"ReviewGrade":  reviewGrade,
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
	effects := map[string]*StatusEffect{}
	for _, skill := range c.Skills {
		for _, effect := range skill.Effects {
			effects[effect.Slug] = effect
		}
	}
	for _, aura := range c.Auras {
		for _, effect := range aura.Effects {
			effects[effect.Slug] = effect
		}
	}
	extraData["Champion"] = c
	extraData["Skills"] = len(c.Skills) + len(c.Auras)
	extraData["Effects"] = effects
	err = tmpl.Execute(output, extraData)
	return err
}

func (c *Champion) SetAura(description string) {
	if len(c.Auras) == 0 {
		c.Auras = append(c.Auras, &Aura{Effects: make([]*StatusEffect, 0)})
	}
	c.Auras[0].RawDescription = description
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

func reviewGrade(gr float64) template.HTML {
	g := ""
	for v, r := range ratingToRank {
		if r == int(gr) {
			g = v
			break
		}
	}
	val := ``
	if g != "" {
		val = fmt.Sprintf(`<span class="champion-rating champion-rating-%s"><strong>%.1f</strong></span> `, g, gr)
	}
	return template.HTML(fmt.Sprintf("%s%s", val, grade(g)))
}

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

func (c *Champion) ParseRawSkill(raw string) error {
	if raw == "" {
		return nil
	} else if strings.Index(raw, "Aura") != -1 {
		return c.setAuraFromRaw(raw)
	}
	return c.setSkillFromRaw(raw)
}

func (c *Champion) setAuraFromRaw(raw string) error {
	data := strings.Join(strings.Split(raw, "\n")[1:], "<br>")
	aura := &Aura{RawDescription: strings.Trim(data, " \n\r")}
	if aura.RawDescription != "" {
		c.Auras = make([]*Aura, 1)
		c.Auras[0] = aura
	}

	return nil
}

func (c *Champion) setSkillFromRaw(raw string) error {
	data := strings.Split(raw, "\n")
	// assuming name is on line 1, before the "Level"
	dataOkForDescription := make([]string, 0)
	for idx := range data {
		data[idx] = strings.Trim(data[idx], " ")
		if idx > 0 && data[idx] != "" {
			fit := false
			concat := false
			switch true {
			case idx == 1,
				strings.HasPrefix(data[idx], "Lvl."),
				strings.HasPrefix(data[idx], "Damage based on:"):
				fit = true
			case data[idx-1] != "" && idx != 1:
				fit = true
				concat = true
			}
			if fit {
				if concat {
					dataOkForDescription[len(dataOkForDescription)-1] = fmt.Sprintf("%s %s", dataOkForDescription[len(dataOkForDescription)-1], data[idx])
				} else {
					dataOkForDescription = append(dataOkForDescription, data[idx])
				}
			}
		}
	}
	text := strings.Join(dataOkForDescription, "<br>")
	name := strings.Trim(data[0][:strings.Index(data[0], "Level")-1], " ")
	okName, errName := GetSanitizedName(name)
	if errName != nil {
		return errName
	}
	for _, skill := range c.Skills {
		if skill.Name == okName {
			skill.RawDescription = text
			return nil
		}
	}
	skill := &Skill{
		Name:           okName,
		RawDescription: text,
	}
	errSanitize := skill.Sanitize()
	if errSanitize != nil {
		return errSanitize
	}
	c.Skills = append(c.Skills, skill)
	return nil
}

func (c *Champion) GetPageExtraData(dataDirectory string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	if dataDirectory == "" {
		return data, nil
	}

	statusList, errStatusEffects := fetchStatusEffects(dataDirectory)
	if errStatusEffects != nil {
		return nil, errStatusEffects
	}

	data["AllEffects"] = statusList

	return data, nil
}

var (
	ErrSkillNotFound = fmt.Errorf("skill not found")
)

func (c *Champion) GetSkillByName(name string) (*Skill, error) {
	for _, skill := range c.Skills {
		if skill.Name == name {
			return skill, nil
		}
	}
	return nil, ErrSkillNotFound
}
