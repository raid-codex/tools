package scrap_raidshadowlegendspro_champions

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionName  *string
	DataDirectory *string
	Skills        *bool
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		DataDirectory: cmd.Flag("data-directory", "Directory containing data").Required().String(),
		ChampionName:  cmd.Flag("champion-name", "Name of the champion being looked up").Required().String(),
		Skills:        cmd.Flag("with-skills", "Fetch champion skills and store them").Bool(),
	}
}

func (c *Command) getFromWebsite(url string) (*http.Response, error) {
	req, errRequest := http.NewRequest("GET", url, nil)
	if errRequest != nil {
		return nil, errRequest
	}
	log.Printf("requesting website... %s\n", req.URL)
	resp, errResponse := http.DefaultClient.Do(req)
	if errResponse != nil {
		return nil, errResponse
	}
	log.Println("got response")
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request %v returned %d", req, resp.StatusCode)
	}
	return resp, nil
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
	resp, errReq := c.getFromWebsite(fmt.Sprintf("https://raidshadowlegends.pro/%s/raid-shadow-legends-%s-build-guide/", champion.FactionSlug, champion.Slug))
	if errReq != nil {
		resp, errReq = c.getFromWebsite(fmt.Sprintf("https://raidshadowlegends.pro/%s/%s/", champion.FactionSlug, champion.Slug))
	}
	if errReq != nil {
		utils.Exit(1, errReq)
	}
	defer resp.Body.Close()
	doc, errDoc := goquery.NewDocumentFromReader(resp.Body)
	if errDoc != nil {
		utils.Exit(1, errDoc)
	}
	if c.Skills != nil && *c.Skills {
		c.parseSkills(champion, doc)
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

var (
	skillNameRpl = regexp.MustCompile("^((.+) Level 1|(Aura))$")
	skipEffects  = regexp.MustCompile("<strong>(.+)</strong>")
	repl         = map[string]string{
		"Decreace":         "Decrease",
		"AIIY":             "Ally",
		"`":                "'",
		"’":                "'",
		"[Passive Effect]": "Passive Effect",
		"[Active Effect]":  "Active Effect",
	}
)

func sanitizeString(str string) string {
	for o, n := range repl {
		str = strings.Replace(str, o, n, -1)
	}
	if strings.HasSuffix(str, "{ ajaxHitsCounterFailedCallback( this );}}}})();") {
		split := strings.Split(str, "<br>")
		str = strings.Join(split[0:len(split)-1], "<br>")
	}
	return str
}

func (c *Command) parseSkills(champion *common.Champion, doc *goquery.Document) {
	doc.Find(".entry-content").Each(func(_ int, s *goquery.Selection) {
		skills := false
		// skils name have h3
		skillNames := map[string]bool{}
		s.Find("h3").Each(func(_ int, sc *goquery.Selection) {
			skillName := strings.Trim(sanitizeString(sc.Text()), " ")
			m := skillNameRpl.FindStringSubmatch(skillName)
			if m == nil {
				utils.Exit(1, fmt.Errorf("Skill name '%s' did not match regex", skillName))
			} else if m[1] == "Aura" {
				m[2] = "Aura"
			}
			skillNames[m[2]] = true
		})
		skillMap := map[string]string{}
		currentSkill := ""
		s.Children().Each(func(_ int, sc *goquery.Selection) {
			skillData := sanitizeString(sc.Text())
			if !skills && strings.HasSuffix(skillData, "Skills") {
				skills = true
			} else if skills {
				if currentSkill == "" {
					for skillName := range skillNames {
						if strings.Contains(skillData, skillName) {
							currentSkill = skillName
							break
						}
					}
				} else {
					html, err := sc.Html()
					if err != nil {
						utils.Exit(1, err)
					}
					html = sanitizeString(html)
					if !skipEffects.MatchString(html) {
						if skillMap[currentSkill] != "" {
							skillMap[currentSkill] += "<br>"
						}
						skillMap[currentSkill] += html
					}
					if strings.Contains(skillData, "Lvl.") || currentSkill == "Aura" {
						// reset
						currentSkill = ""
					}
				}
			}
		})
		for skillName, skillDescription := range skillMap {
			if skillName == "Aura" {
				champion.Auras = []*common.Aura{
					&common.Aura{
						RawDescription: strings.ToUpper(skillDescription[0:1]) + skillDescription[1:],
					},
				}
			} else {
				passive := false
				if strings.Contains(skillName, "[P]") {
					skillName = strings.Replace(skillName, "[P]", "", -1)
					passive = true
				}
				skillName = strings.Trim(skillName, " ")
				skill, errSkill := champion.GetSkillByName(skillName)
				if errSkill != nil {
					// try with invalid name that we could have pushed
					skill, errSkill = champion.GetSkillByName(skillName + " Level 1")
				}
				if errSkill != nil {
					skill = &common.Skill{Name: skillName}
					champion.Skills = append(champion.Skills, skill)
				}
				if strings.Contains(skillName, "[P]") {
					skillName = strings.Trim(strings.Replace(skillName, "[P]", "", -1), " ")
					passive = true
				}
				skill.Passive = passive
				skill.Name = skillName
				if skillDescription != "" {
					skill.RawDescription = skillDescription
				}
			}
		}
	})
}
