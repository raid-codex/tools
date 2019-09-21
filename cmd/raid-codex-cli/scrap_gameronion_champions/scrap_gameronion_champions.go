package scrap_gameronion_champions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-test/deep"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionName  *string
	DataDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		DataDirectory: cmd.Flag("data-directory", "Directory containing data").Required().String(),
		ChampionName:  cmd.Flag("champion-name", "Name of the champion being looked up").Required().String(),
	}
}

func (c *Command) Run() {
	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	champions, errChampions := common.GetChampions(common.FilterChampionName(*c.ChampionName))
	if errChampions != nil {
		utils.Exit(1, errChampions)
	} else if len(champions) != 1 {
		utils.Exit(1, fmt.Errorf("found %d champions with name %s", len(champions), *c.ChampionName))
	}
	champion := champions[0]
	req, errRequest := http.NewRequest("GET", fmt.Sprintf("https://www.gameronion.com/Raid-Shadow-Legends/champions/%s", champion.Slug), nil)
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
	doc.Find(".faction p").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			faction := strings.Trim(s.Find("a").Text(), " \n")
			log.Printf("Faction: %s\n", faction)
		}
	})
	doc.Find(".rarity").Each(func(i int, s *goquery.Selection) {
		log.Printf("Rarity: %s\n", strings.Trim(s.Text(), " \n"))
	})
	doc.Find(".champtype").Each(func(i int, s *goquery.Selection) {
		log.Printf("Type: %s\n", strings.Trim(s.Text(), " \n"))
	})
	doc.Find(".skill-cont").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		// remove double \n + double spaces
		txt = emptyLinesRegexp.ReplaceAllString(txt, "")
		for _, repl := range replacements {
			for strings.Index(txt, repl[0]) != -1 {
				txt = strings.Replace(txt, repl[0], repl[1], -1)
			}
		}
		var skillName, skillDescription string
		var cooldown int64
		var passive bool
		parts := strings.Split(txt, "\n")
		for _, part := range parts {
			part = strings.Trim(part, " ")
			if part == "" {
				continue
			}
			switch true {
			case strings.HasPrefix(part, "[Passive]"):
				passive = true
			case strings.HasPrefix(part, "Level "):
				continue
			case strings.HasPrefix(part, "Cooldown: "):
				pCooldown, errCooldown := strconv.ParseInt(strings.Trim(part[9:], " "), 10, 64)
				if errCooldown != nil {
					utils.Exit(1, errCooldown)
				}
				cooldown = pCooldown
			case strings.HasPrefix(part, "Lvl"):
				skillDescription = fmt.Sprintf("%s<br>%s", skillDescription, part)
			default:
				if skillName == "" {
					skillName = part
					if strings.Contains(skillName, "Aura") {
						skillName = "Aura"
					}
				} else {
					skillDescription = part
				}
			}
		}
		if skillName != "Aura" {
			//log.Printf("Skill: %s\nPassive: %t\nCooldown: %d\nDescription:\n\t%s\n", skillName, passive, cooldown, strings.Replace(skillDescription, "<br>", "\n\t", -1))
			var before []byte
			skill, errSkill := champion.GetSkillByName(skillName)
			if errSkill != nil {
				skill = &common.Skill{Name: skillName}
				champion.Skills = append(champion.Skills, skill)
			} else {
				before, _ = json.Marshal(skill)
			}
			if skillDescription != "" {
				skill.RawDescription = skillDescription
			}
			if passive {
				skill.Passive = true
			}
			if len(skill.Upgrades) == 0 {
				skill.Upgrades = append(skill.Upgrades, &common.SkillData{
					Target: &common.Target{},
				})
			}
			for _, upgrade := range skill.Upgrades {
				if upgrade.Cooldown == 0 {
					upgrade.Cooldown = cooldown
				}
			}
			after, _ := json.Marshal(skill)
			log.Printf("Diff for <%s>\n", skill.Name)
			printDiff(before, after)
		} else {
			//log.Printf("Aura:\n\t%s\n", skillDescription)
			champion.Auras = []*common.Aura{
				&common.Aura{
					RawDescription: skillDescription,
				},
			}
		}
	})
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

var (
	replacements = [3][2]string{
		{"  ", " "},
		{" \n", "\n"},
		{"\n\n", "\n"},
	}
	emptyLinesRegexp = regexp.MustCompile(`(?m)^(?:[\t ]*(?:\r?\n|\r))+`)
)

func printDiff(a []byte, b []byte) {
	aa := map[string]interface{}{}
	bb := map[string]interface{}{}
	_ = json.Unmarshal(a, &aa)
	_ = json.Unmarshal(b, &bb)
	if diff := deep.Equal(aa, bb); diff != nil {
		log.Printf("Diff: %v\n", diff)
	} else {
		log.Println("no diff")
	}
}
