package scrap_ayumilove_champions

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionSlug  *string
	DataDirectory *string
	Stats         *bool
	Builds        *bool
	Masteries     *bool
	Ratings       *bool
	Skills        *bool
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		DataDirectory: cmd.Flag("data-directory", "Directory containing data").Required().String(),
		ChampionSlug:  cmd.Flag("champion-slug", "Slug of the champion being looked up").Required().String(),
		Stats:         cmd.Flag("with-stats", "Fetch champion stats and store them").Bool(),
		Builds:        cmd.Flag("with-builds", "Fetch and store champion's build").Bool(),
		Masteries:     cmd.Flag("with-masteries", "Fetch and store champion's masteries").Bool(),
		Ratings:       cmd.Flag("with-ratings", "Fetch and store champion's rating").Bool(),
		Skills:        cmd.Flag("with-skills", "Also parse champions' skills").Bool(),
	}
}

func (c *Command) Run() {
	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	champions, errChampions := common.GetChampions(common.FilterChampionSlug(*c.ChampionSlug))
	if errChampions != nil {
		utils.Exit(1, errChampions)
	} else if len(champions) != 1 {
		utils.Exit(1, fmt.Errorf("found %d champions with slug %s", len(champions), *c.ChampionSlug))
	}
	champion := champions[0]
	slugToLookup := champion.Slug
	switch slugToLookup {
	case "ma-shalled":
		slugToLookup = "mashalled"
	}
	req, errRequest := http.NewRequest("GET", fmt.Sprintf("https://ayumilove.net/raid-shadow-legends-%s-skill-mastery-equip-guide/", slugToLookup), nil)
	if errRequest != nil {
		utils.Exit(1, errRequest)
	}
	resp, errResponse := http.DefaultClient.Do(req)
	if errResponse != nil {
		utils.Exit(1, errResponse)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		utils.Exit(1, fmt.Errorf("request %v returned %d", req, resp.StatusCode))
	}
	doc, errDoc := goquery.NewDocumentFromReader(resp.Body)
	if errDoc != nil {
		utils.Exit(1, errDoc)
	}
	if c.Builds != nil && *c.Builds {
		// don't keep ayumilove's builds
		builds := []*common.Build{}
		for _, build := range champion.RecommendedBuilds {
			if build.From != "ayumilove.net" {
				builds = append(builds, build)
			}
		}
		champion.RecommendedBuilds = builds
		c.parseEquipment(champion, doc)
	}
	if c.Masteries != nil && *c.Masteries {
		// don't keep ayumilove's masteries
		masteries := []*common.ChampionMasteries{}
		for _, mastery := range champion.Masteries {
			if mastery.From != "ayumilove.net" {
				masteries = append(masteries, mastery)
			}
		}
		champion.Masteries = masteries
		c.parseMasteries(champion, doc)
	}
	if c.Stats != nil && *c.Stats {
		c.parseStats(champion, doc)
	}
	if c.Ratings != nil && *c.Ratings {
		c.parseRating(champion, doc)
	}
	if c.Skills != nil && *c.Skills {
		//c.parseSkills(champion, doc)
	}
	errSanitize := champion.Sanitize()
	if errSanitize != nil {
		utils.Exit(1, errSanitize)
	}
	// write champion file
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/docs/champions/current/%s.json", *c.DataDirectory, champion.Slug), champion)
	if errWrite != nil {
		utils.Exit(1, errWrite)
	}
}

func (c *Command) parseStats(champion *common.Champion, doc *goquery.Document) {
	doc.Find(".entry-content table").Each(func(idx int, s *goquery.Selection) {
		if idx != 0 {
			// only the first index is interesting for stats
			return
		}
		s.Find("td").Each(func(subIdx int, sc *goquery.Selection) {
			if subIdx != 2 {
				// only the 3rd column is for stats
				return
			}
			data := strings.Split(sc.Text(), "\n")
			chars := champion.Characteristics[60]
			for _, d := range data {
				if len(d) == 0 || strings.Contains(d, "Total Stats") {
					continue
				}
				var intField *int64
				var floatField *float64
				switch true {
				case strings.Contains(d, "ATK"):
					intField = &chars.Attack
				case strings.Contains(d, "HP"):
					intField = &chars.HP
				case strings.Contains(d, "DEF"):
					intField = &chars.Defense
				case strings.Contains(d, "ACC"):
					intField = &chars.Accuracy
				case strings.Contains(d, "RESIST"):
					intField = &chars.Resistance
				case strings.Contains(d, "SPD"):
					intField = &chars.Speed
				case strings.Contains(d, "C. Rate"):
					floatField = &chars.CriticalRate
				case strings.Contains(d, "C. DMG"):
					floatField = &chars.CriticalDamage
				default:
					utils.Exit(1, fmt.Errorf("cannot parse stats line '%s'", d))
				}
				subD := strings.Split(d, " ")
				lastPart := strings.Replace(strings.Replace(subD[len(subD)-1], "%", "", -1), ",", "", -1)
				v, errInt := strconv.ParseInt(lastPart, 10, 64)
				if errInt != nil {
					utils.Exit(1, fmt.Errorf("invalid number: %s ; %s", lastPart, errInt))
				}
				if intField != nil {
					*intField = v
				} else if floatField != nil {
					*floatField = float64(v) / 100.0
				}
			}
			champion.Characteristics[60] = chars
		})
	})
}

var (
	rankRegexp    = regexp.MustCompile(`★`)
	countToLetter = map[int]string{
		5: "SS",
		4: "S",
		3: "A",
		2: "B",
		1: "C",
		0: "D",
	}
)

func (c *Command) parseRating(champion *common.Champion, doc *goquery.Document) {
	rating := &common.Rating{}
	doc.Find(".entry-content table").Each(func(idx int, s *goquery.Selection) {
		if idx != 1 {
			// only the second index is interesting for stats
			return
		}
		s.Find("td").Each(func(subIdx int, sc *goquery.Selection) {
			data := strings.Split(sc.Text(), "\n")
			for _, d := range data {
				count := len(rankRegexp.FindAllStringIndex(d, -1))
				rank := countToLetter[count]
				switch true {
				case strings.Contains(d, "Campaign"):
					rating.Campaign = rank
				case strings.Contains(d, "Arena Defense"):
					rating.ArenaDef = rank
				case strings.Contains(d, "Arena Offense"):
					rating.ArenaOff = rank
				case strings.Contains(d, "Clan Boss"):
					rating.ClanBossWoGS = rank
					rating.ClanBosswGS = rank
				case strings.Contains(d, "Minotaur"):
					rating.Minotaur = rank
				case strings.Contains(d, "Spider"):
					rating.Spider = rank
				case strings.Contains(d, "Fire Knight"):
					rating.FireKnight = rank
				case strings.Contains(d, "Dragon"):
					rating.Dragon = rank
				case strings.Contains(d, "Ice Golem"):
					rating.IceGuardian = rank
				case strings.Contains(d, "Void Keep"):
					rating.VoidDungeon = rank
				case strings.Contains(d, "Magic Keep"):
					rating.MagicDungeon = rank
				case strings.Contains(d, "Force Keep"):
					rating.ForceDungeon = rank
				case strings.Contains(d, "Spirit Keep"):
					rating.SpiritDungeon = rank
				}

			}
		})
	})
	champion.AddRating("ayumilove", rating, 5)
}

var (
	equipmentPrefixes = []string{"Weapon", "Helmet", "Shield", "Gauntlets", "Chestplate", "Boots", "Ring", "Amulet", "Banner"}
	mainStatExtracter = regexp.MustCompile(`\((.+)\)$`)
	knownLocations    = []string{"Arena", "Campaign", "Clan Boss", "Dungeon", "Faction Wars"}
	setExtracter      = regexp.MustCompile(`(\d) ([A-Za-z ]+) Set`)
)

func (c *Command) parseEquipment(champion *common.Champion, doc *goquery.Document) {
	equipmentContent := []string{}
	doc.Find(".entry-content").Each(func(_ int, s *goquery.Selection) {
		check := 0
		s.Children().Each(func(_ int, sc *goquery.Selection) {
			switch check {
			case 1:
				if sc.Is("h2") || strings.HasSuffix(sc.Text(), "Mastery Guide") {
					check = 2
				} else {
					equipmentContent = append(equipmentContent, sc.Text())
				}
			case 2:
				// mastery guide
				if sc.Is("table") {
					check = 0
				}
			default:
				check = 0
				if strings.HasSuffix(sc.Text(), "Equipment Guide") {
					check = 1
				} else if strings.HasSuffix(sc.Text(), "Mastery Guide") {
					check = 2
				}
			}
		})
	})
	parseEquipment(champion, strings.Join(equipmentContent, "\n"))
}

func (c *Command) parseMasteries(champion *common.Champion, doc *goquery.Document) {
	content := []string{}
	doc.Find(".entry-content").Each(func(_ int, s *goquery.Selection) {
		check := 0
		s.Children().Each(func(_ int, sc *goquery.Selection) {
			switch check {
			case 2:
				// mastery guide
				if sc.Is("table") {
					sc.Find("tr td ol li").Each(func(_ int, mastery *goquery.Selection) {
						content = append(content, fmt.Sprintf("mastery: %s", mastery.Text()))
					})
				} else if sc.Is("h2") {
					check = 0
					break
				} else {
					content = append(content, sc.Text())
				}
			default:
				check = 0
				if strings.HasSuffix(sc.Text(), "Mastery Guide") {
					check = 2
				}
			}
		})
	})
	parseMasteries(champion, strings.Join(content, "\n"))
}

func parseMasteries(champion *common.Champion, content string) {
	chunks := strings.Split(content, "\n")
	masteries := []*common.ChampionMasteries{}
	var currentMastery *common.ChampionMasteries
	idx := 0
	for idx < len(chunks) {
		chunk := chunks[idx]
		if strings.HasPrefix(chunk, "mastery: ") {
			if currentMastery == nil {
				// assume everywhere
				currentMastery = &common.ChampionMasteries{
					From:      "ayumilove.net",
					Author:    "ayumilove",
					Locations: []string{},
				}
				masteries = append(masteries, currentMastery)
			}
			mastery := chunk[9:]
			if mastery != "N/A" {
				mastery = knownMasteriesReplacement(mastery)
				found, err := common.GetMasteries(common.FilterMasteryLowercasedName(mastery))
				if err != nil {
					utils.Exit(1, fmt.Errorf("mastery %s not found", mastery))
				} else if len(found) != 1 {
					panic(fmt.Errorf("mastery %s found %d times", mastery, len(found)))
					//utils.Exit(1, fmt.Errorf("mastery %s found %d times", mastery, len(found)))
				}
				switch found[0].Tree {
				case 1:
					currentMastery.Offense = append(currentMastery.Offense, found[0].Slug)
				case 2:
					currentMastery.Defense = append(currentMastery.Defense, found[0].Slug)
				case 3:
					currentMastery.Support = append(currentMastery.Support, found[0].Slug)
				default:
					utils.Exit(1, fmt.Errorf("invalid tree %d", found[0].Tree))
				}
			}
		} else if len(chunk) > 0 {
			// where
			currentMastery = &common.ChampionMasteries{
				From:      "ayumilove.net",
				Author:    "ayumilove",
				Locations: []string{},
			}
			for _, location := range knownLocations {
				if strings.Contains(chunk, location) {
					currentMastery.Locations = append(currentMastery.Locations, common.GetLinkNameFromSanitizedName(location))
				}
			}
			masteries = append(masteries, currentMastery)
		}
		idx++
	}
	for _, mastery := range masteries {
		champion.AddMastery(mastery)
	}
}

var (
	knownMasteriesErrors = map[string]string{
		"Swam Smiter":              "Swarm Smiter",
		"Eagle-Eye":                "Eagle Eye",
		"Blood Thirst":             "Bloodthirst",
		"Whirldwind of Death":      "Whirlwind of Death",
		"Subborness":               "Stubborness",
		"Pintpoint Accuracy":       "Pinpoint Accuracy",
		"Shiedl Breaker":           "Shield Breaker",
		"Stubbornness":             "Stubborness",
		"Delay of Death":           "Delay Death",
		"Warmaster / Giant Slayer": "Warmaster",
	}
)

func knownMasteriesReplacement(mastery string) string {
	if v, ok := knownMasteriesErrors[mastery]; ok {
		return v
	}
	return mastery
}

func parseEquipment(champion *common.Champion, equipment string) {
	builds := []*common.Build{}
	chunks := strings.Split(equipment, "\n")
	statPrio := []string{}
	idx := 0
	for idx < len(chunks) {
		chunk := chunks[idx]
		if strings.HasPrefix(chunk, "Equipment Set for") {
			chunks[idx] = strings.Replace(chunks[idx], "Equipment Set for Campaign, Clan Boss, Dungeon, 1", "Equipment Set for Campaign, Clan Boss, Dungeon: 1", 1)
			chunk = chunks[idx]
			if strings.HasPrefix(chunks[idx], "Equipment Set for") && strings.Contains(chunks[idx], ":") {
				for strings.HasPrefix(chunks[idx], "Equipment Set for") && strings.Contains(chunks[idx], ":") {
					builds = append(builds, parseSet(strings.Split(chunks[idx], ":")[1], parseLocations(chunks[idx])))
					idx++
				}
				idx--
			} else {
				locations := parseLocations(chunk)
				idx++
				for !strings.HasPrefix(chunks[idx], "Equipment") {
					builds = append(builds, parseSet(chunks[idx], locations))
					idx++
				}
				idx--
			}
		} else if chunk == "Equipment Stat Priority" {
			idx++
			statPrio = parseStatPrio(chunks[idx])
		} else if strings.HasPrefix(chunk, "Stat Priority: ") {
			statPrio = parseStatPrio(chunk[14:])
		} else {
			for _, prefix := range equipmentPrefixes {
				if strings.HasPrefix(chunk, prefix) {
					mainStatExtract := mainStatExtracter.FindStringSubmatch(chunk)
					for _, build := range builds {
						sp := &common.StatPriority{
							MainStat:        mainStatExtract[1],
							AdditionalStats: statPrio[:],
						}
						build.Set(prefix, sp)
					}
				}
			}
		}
		idx++
	}
	for _, build := range builds {
		errSanitize := build.Sanitize()
		if errSanitize != nil {
			utils.Exit(1, errSanitize)
		}
		champion.AddBuild(build)
	}
}

func parseLocations(chunk string) []string {
	locations := []string{}
	for _, location := range knownLocations {
		if strings.Contains(chunk, location) {
			locations = append(locations, common.GetLinkNameFromSanitizedName(location))
		}
	}
	return locations
}

func parseSet(chunk string, locations []string) *common.Build {
	setsExtract := setExtracter.FindAllStringSubmatch(chunk, -1)
	sets := []string{}
	for _, setExtract := range setsExtract {
		nbr, _ := strconv.Atoi(setExtract[1])
		idx2 := 0
		for idx2 < nbr {
			idx2++
			sets = append(sets, common.GetLinkNameFromSanitizedName(setExtract[2]))
		}
	}
	return &common.Build{
		From:      "ayumilove.net",
		Author:    "ayumilove",
		Locations: locations,
		Sets:      sets,
	}
}

func parseStatPrio(chunk string) []string {
	statPrio := []string{}
	miniChunks := strings.Split(chunk, ",")
	for _, stat := range miniChunks {
		statPrio = append(statPrio, strings.Trim(stat, " "))
	}
	return statPrio
}
