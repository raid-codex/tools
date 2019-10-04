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
	req, errRequest := http.NewRequest("GET", fmt.Sprintf("https://ayumilove.net/raid-shadow-legends-%s-skill-mastery-equip-guide/", champion.Slug), nil)
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
	doc.Find(".entry-content").Each(func(_ int, s *goquery.Selection) {
		check := 0
		s.Children().Each(func(_ int, sc *goquery.Selection) {
			switch check {
			case 1:
				// equipment guide
				builds := parseEquipment(sc, champion)
				for _, build := range builds {
					champion.AddBuild(build)
				}
				check = 0
			case 2:
				// mastery guide
				if sc.Is("table") {
					sc.Find("tr td ol li").Each(func(_ int, mastery *goquery.Selection) {
						//log.Printf("%v\n", mastery.Text())
					})
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
	equipmentPrefixes = []string{"Weapon", "Helmet", "Shield", "Gauntlets", "Chestplate", "Boots", "Ring", "Amulet", "Banner"}
	mainStatExtracter = regexp.MustCompile(`\((.+)\)$`)
	knownLocations    = []string{"Arena", "Campaign", "Clan Boss", "Dungeon"}
	setExtracter      = regexp.MustCompile(`(\d) ([A-Za-z ]+) Set`)
)

func parseEquipment(sc *goquery.Selection, champion *common.Champion) []*common.Build {
	builds := []*common.Build{}
	chunks := strings.Split(sc.Text(), "\n")
	statPrio := []string{}
	for _, chunk := range chunks {
		if strings.HasPrefix(chunk, "Equipment Set for") {
			locations := []string{}
			for _, location := range knownLocations {
				if strings.Contains(chunk, location) {
					locations = append(locations, common.GetLinkNameFromSanitizedName(location))
				}
			}
			setsExtract := setExtracter.FindAllStringSubmatch(chunk, -1)
			sets := []string{}
			for _, setExtract := range setsExtract {
				nbr, _ := strconv.Atoi(setExtract[1])
				idx := 0
				for idx < nbr {
					idx++
					sets = append(sets, common.GetLinkNameFromSanitizedName(setExtract[2]))
				}
			}
			builds = append(builds, &common.Build{
				From:      "ayumilove.net",
				Author:    "ayumilove",
				Locations: locations,
				Sets:      sets,
			})
		} else if strings.HasPrefix(chunk, "Stat Priority: ") {
			chunk = strings.Replace(chunk, "Stat Priority: ", "", -1)
			miniChunks := strings.Split(chunk, ",")
			for _, stat := range miniChunks {
				statPrio = append(statPrio, strings.Trim(stat, " "))
			}
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
	}
	for _, build := range builds {
		errSanitize := build.Sanitize()
		if errSanitize != nil {
			utils.Exit(1, errSanitize)
		}
	}
	return builds
}
