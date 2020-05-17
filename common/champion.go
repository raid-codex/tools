package common

import (
	"crypto/md5"
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
	DateAdded          string                    `json:"date_added"`
	Rarity             string                    `json:"rarity"`
	Element            string                    `json:"element"`
	Type               string                    `json:"type"`
	Rating             *Rating                   `json:"rating"`
	Reviews            *Review                   `json:"reviews"`
	AllRatings         AllRatings                `json:"all_ratings"`
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
	Synergies          []*Synergy                `json:"synergy"`
	Thumbnail          string                    `json:"thumbnail"`
	Tags               []string                  `json:"tags"`
	Masteries          []*ChampionMasteries      `json:"masteries"`
	FusionData         []*ChampionFusionData     `json:"fusion_data"`
	EffectSlugs        []string                  `json:"effect_slugs"`
	Videos             []*Video                  `json:"videos"`
}

type ChampionFusionData struct {
	FusionSlug string `json:"fusion_slug"`
	FusionType string `json:"fusion_type"`
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
	for _, build := range c.RecommendedBuilds {
		errSanitize := build.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
	}

	c.Slug = c.LinkName()

	// faction
	// -- for compat
	if c.FactionSlug == "skinwalker" {
		c.FactionSlug = "skinwalkers"
	}
	factions, errFactions := GetFactions(FilterFactionSlug(c.FactionSlug))
	if errFactions != nil {
		return errFactions
	} else if len(factions) != 1 {
		return fmt.Errorf("found %d factions with slug %s", len(factions), c.FactionSlug)
	}
	faction := factions[0]
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

	passiveCount := 1
	skillCount := 1
	for idx, skill := range c.Skills {
		if c.GIID != "" && skill.GIID == "" {
			skill.GIID = fmt.Sprintf("%s_s%d", c.GIID, idx+1)
		}
		if skill.SkillNumber == "" {
			if skill.Passive {
				skill.SkillNumber = fmt.Sprintf("P%d", passiveCount)
			} else {
				skill.SkillNumber = fmt.Sprintf("A%d", skillCount)
			}
		}
		errSanitize := skill.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
		if skill.Passive {
			passiveCount++
		} else {
			skillCount++
		}
	}

	for _, aura := range c.Auras {
		errSanitize := aura.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
	}

	if c.Synergies == nil {
		c.Synergies = make([]*Synergy, 0)
	}
	errSynergy := c.computeSynergy()
	if errSynergy != nil {
		return errSynergy
	}
	for _, synergy := range c.Synergies {
		errSanitizeSynergy := synergy.Sanitize()
		if errSanitizeSynergy != nil {
			return errSanitizeSynergy
		}
	}

	if c.GIID != "" && c.Thumbnail != "unknown" {
		c.Thumbnail = fmt.Sprintf("%x", md5.Sum([]byte(c.GIID)))
	} else {
		c.Thumbnail = "unknown"
	}

	if c.Tags == nil {
		c.Tags = make([]string, 0)
	}
	sort.SliceStable(c.Tags, func(i, j int) bool { return c.Tags[i] < c.Tags[j] })

	if c.Masteries == nil {
		c.Masteries = make([]*ChampionMasteries, 0)
	}
	for _, mastery := range c.Masteries {
		if err := mastery.Sanitize(); err != nil {
			return err
		}
	}

	if err := c.lookupFusions(); err != nil {
		return err
	}

	effectSlugs := map[string]bool{}
	for _, skill := range c.Skills {
		for _, effect := range skill.Effects {
			effectSlugs[effect.Slug] = true
		}
		for _, upgrade := range skill.Upgrades {
			for _, effect := range upgrade.Effects {
				effectSlugs[effect.Slug] = true
			}
		}
	}
	c.EffectSlugs = make([]string, len(effectSlugs))
	idx := 0
	for slug := range effectSlugs {
		c.EffectSlugs[idx] = slug
		idx++
	}
	sort.SliceStable(c.EffectSlugs, func(i, j int) bool { return c.EffectSlugs[i] < c.EffectSlugs[j] })

	if err := c.sanitizeVideos(); err != nil {
		return err
	}

	if c.AllRatings == nil {
		c.AllRatings = make(AllRatings, 0)
	}
	for _, r := range c.AllRatings {
		if err := r.Sanitize(); err != nil {
			return err
		}
	}

	c.Rating = c.AllRatings.Compute()
	if err := c.Rating.Sanitize(); err != nil {
		return err
	}

	// don't store champion slugs on factions / effects
	c.Faction.ChampionSlugs = []string{}
	for _, skill := range c.Skills {
		for _, effect := range skill.Effects {
			effect.ChampionSlugs = []string{}
		}
		for _, upgrade := range skill.Upgrades {
			for _, effect := range upgrade.Effects {
				effect.ChampionSlugs = []string{}
			}
		}
	}

	return nil
}

func (c *Champion) sanitizeVideos() error {
	if c.Videos == nil {
		c.Videos = make([]*Video, 0)
	}
	videos := make(map[string]*Video)
	for _, video := range c.Videos {
		if err := video.Sanitize(); err != nil {
			return err
		}
		key := fmt.Sprintf("%s-%s", video.Source, video.ID)
		if _, ok := videos[key]; ok {
			if videos[key].DateAdded > video.DateAdded {
				videos[key] = video
			}
		} else {
			videos[key] = video
		}
	}
	c.Videos = make([]*Video, len(videos))
	idx := 0
	for _, video := range videos {
		c.Videos[idx] = video
		idx++
	}
	sort.SliceStable(c.Videos, func(i, j int) bool {
		if c.Videos[i].Source == c.Videos[j].Source {
			if c.Videos[i].DateAdded == c.Videos[j].DateAdded {
				return c.Videos[i].ID < c.Videos[j].ID
			}
			return c.Videos[i].DateAdded < c.Videos[j].DateAdded
		}
		return c.Videos[i].Source < c.Videos[j].Source
	})
	return nil
}

func (c *Champion) AddRating(source string, rating *Rating, weight int) {
	for _, src := range c.AllRatings {
		if src.Source == source {
			src.Rating = rating
			src.Weight = weight
			return
		}
	}
	c.AllRatings = append(c.AllRatings, &RatingSource{
		Source: source,
		Rating: rating,
		Weight: weight,
	})
}

func (c *Champion) lookupFusions() error {
	fusions, errFusions := GetFusions(func(f *Fusion) bool {
		if f.hasChampion(c.Slug) {
			return true
		}
		return false
	})
	if errFusions != nil {
		return errFusions
	}
	c.FusionData = make([]*ChampionFusionData, 0)
	for _, fusion := range fusions {
		if fusion.ChampionSlug == c.Slug {
			c.FusionData = append(c.FusionData, &ChampionFusionData{FusionSlug: fusion.Slug, FusionType: "fused"})
		} else if fusion.hasChampionAsIngredient(c.Slug) {
			c.FusionData = append(c.FusionData, &ChampionFusionData{FusionSlug: fusion.Slug, FusionType: "ingredient"})
		}
	}
	return nil
}

func (c *Champion) defaultRating() {
	if c.Rating == nil {
		c.Rating = &Rating{}
	}
	if c.Reviews == nil {
		c.Reviews = &Review{}
	}
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
	name = strings.Trim(name, `!`)
	for _, part := range []string{
		`'`, `"`, ` `, `_`, `(P)`, `+`, `%`, `!`, `’`,
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

func (c *Champion) GetPageContent_Templates(tmpl *template.Template, output io.Writer, extraData map[string]interface{}) error {
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
	return tmpl.Execute(output, extraData)
}

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

func (cl ChampionList) Union(oth ChampionList) ChampionList {
	unique := map[string]*Champion{}
	for _, c := range cl {
		unique[c.Slug] = c
	}
	for _, c := range oth {
		unique[c.Slug] = c
	}
	newList := make(ChampionList, len(unique))
	idx := 0
	for _, champion := range unique {
		newList[idx] = champion
		idx++
	}
	newList.Sort()
	return newList
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

	fusions, errFusions := GetFusions()
	if errFusions != nil {
		return nil, errFusions
	}
	fusionsM := map[string]*Fusion{}
	for _, fusion := range fusions {
		fusionsM[fusion.Slug] = fusion
	}
	data["Fusions"] = fusionsM

	champions, errChampions := GetChampions()
	if errChampions != nil {
		return nil, errChampions
	}
	championsM := map[string]*Champion{}
	for _, champion := range champions {
		championsM[champion.Slug] = champion
	}
	data["Champions"] = championsM

	return data, nil
}

var (
	ErrSkillNotFound = fmt.Errorf("skill not found")
)

func (c *Champion) AddSkill(name string, description string, passive bool) *Skill {
	skill, err := c.GetSkillByName(name)
	if err != nil {
		skill, err = c.GetSkillBySlug(GetLinkNameFromSanitizedName(name))
		if err != nil {
			skill = &Skill{}
			c.Skills = append(c.Skills, skill)
		}
	}
	skill.Name = name
	skill.RawDescription = description
	skill.Passive = passive
	return skill
}

func (c *Champion) GetSkillBySlug(slug string) (*Skill, error) {
	for _, skill := range c.Skills {
		if skill.Slug == slug {
			return skill, nil
		}
	}
	return nil, ErrSkillNotFound
}

func (c *Champion) GetSkillByName(name string) (*Skill, error) {
	for _, skill := range c.Skills {
		if skill.Name == name {
			return skill, nil
		}
	}
	return nil, ErrSkillNotFound
}

func (c *Champion) computeSynergy() error {
	if err := c.synergyA1Poison(); err != nil {
		return err
	}
	if err := c.synergyCounterAttack(); err != nil {
		return err
	}
	return nil
}

var (
	allyCounterattack = FilterChampionStatusEffectWithTargets(StatusEffect_CounterAttack, TargetWho_AllAlly, TargetWho_TargetAlly, TargetWho_OtherAlly)
)

func (c *Champion) synergyA1Poison() error {
	switch true {
	case FilterChampionStatusEffectOnSkill("A1", "poison")(c), FilterChampionStatusEffectOnSkill("A1", "poison-2")(c):
		break
	default:
		return nil
	}
	// A1 poison is good with counterattack
	counterAttack, errListCounterattack := GetChampions(allyCounterattack, FilterChampionNotSlug(c.Slug))
	if errListCounterattack != nil {
		return errListCounterattack
	}
	if len(counterAttack) > 0 {
		synergy := c.getSynergy(SynergyContextKey_PoisonCounterattack)
		synergy.Champions = make([]string, len(counterAttack))
		for idx, champion := range counterAttack {
			synergy.Champions[idx] = champion.Slug
		}
	}
	return nil
}

func (c *Champion) synergyCounterAttack() error {
	if !allyCounterattack(c) {
		return nil
	}
	// Look for A1 poison
	championsPoison1, errPoison1 := GetChampions(FilterChampionStatusEffectOnSkill("A1", "poison"), FilterChampionNotSlug(c.Slug))
	if errPoison1 != nil {
		return errPoison1
	}
	championsPoison2, errPoison2 := GetChampions(FilterChampionStatusEffectOnSkill("A1", "poison-2"), FilterChampionNotSlug(c.Slug))
	if errPoison2 != nil {
		return errPoison2
	}
	championsPoison := championsPoison1.Union(championsPoison2)
	if len(championsPoison) > 0 {
		synergy := c.getSynergy(SynergyContextKey_PoisonCounterattack)
		synergy.Champions = make([]string, len(championsPoison))
		for idx, champion := range championsPoison {
			synergy.Champions[idx] = champion.Slug
		}
	}
	return nil
}

func (c *Champion) getSynergy(key SynergyContextKey) *Synergy {
	for _, synergy := range c.Synergies {
		if synergy.Context.Key == key {
			return synergy
		}
	}
	synergy := &Synergy{Context: SynergyContext{Key: key}}
	c.Synergies = append(c.Synergies, synergy)
	return synergy
}

func (c *Champion) AddBuild(build *Build) {
	c.RecommendedBuilds = append(c.RecommendedBuilds, build)
}

func (c *Champion) AddMastery(mastery *ChampionMasteries) {
	c.Masteries = append(c.Masteries, mastery)
}
