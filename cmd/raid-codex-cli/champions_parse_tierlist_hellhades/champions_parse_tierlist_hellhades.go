package champions_parse_tierlist_hellhades

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

//

type Command struct {
	CSVFile       *string
	DataDirectory *string
}

var (
	ErrNotFound = fmt.Errorf("not found")
	nameReplace = map[string]string{
		"Centurian":             "Centurion",
		"Steadfast Marshall":    "Steadfast Marshal",
		"Woad Painted":          "Woad-Painted",
		"Ma’Shalled":            "Ma'Shalled",
		"Big ‘Un":               "Big'Un",
		"Nogdar The Headhunter": "Nodgar the Headhunter",
	}
)

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		DataDirectory: cmd.Flag("data-directory", "Data directory").Required().String(),
	}
	return command
}

const safeguard = `NameOverall RatingClan BossFaction WarsSpiderDragonFire KnightIce GolemArena DefArena Atk`

func (c *Command) Run() {
	errInit := common.InitFactory(*c.DataDirectory)
	if errInit != nil {
		utils.Exit(1, errInit)
	}
	doc, errDoc := c.requestUrl("https://www.hellhades.com/raid-shadow-legends-tier-list/")
	if errDoc != nil {
		utils.Exit(1, errDoc)
	}
	errors := make([]error, 0)
	champions := make([]*common.Champion, 0)
	doc.Find(".post-content table.posts-data-table tr").Each(func(idx int, s *goquery.Selection) {
		if idx == 0 {
			if s.Text() != safeguard {
				utils.Exit(1, fmt.Errorf("invalid safe guard: '%s' instead of '%s'", s.Text(), safeguard))
			}
			return
		}
		var champion *common.Champion
		var rating common.Rating
		var name string
		s.Find("td").Each(func(cIdx int, cS *goquery.Selection) {
			switch cIdx {
			case col_Name:
				name = cS.Text()
				if nameReplace[name] != "" {
					name = nameReplace[name]
				}
				champions, errChampion := common.GetChampions(func(c *common.Champion) bool {
					return strings.ToLower(c.Name) == strings.ToLower(name)
				})
				if errChampion != nil {
					utils.Exit(1, errChampion)
				} else if len(champions) == 1 {
					champion = champions[0]
				}
			case col_ClanBoss:
				rating.ClanBossWoGS = sanitizeRating(cS.Text())
				rating.ClanBosswGS = sanitizeRating(cS.Text())
			case col_FactionWars:
				rating.FactionWars = sanitizeRating(cS.Text())
			case col_Spider:
				rating.Spider = sanitizeRating(cS.Text())
			case col_Dragon:
				rating.Dragon = sanitizeRating(cS.Text())
			case col_FireKnight:
				rating.FireKnight = sanitizeRating(cS.Text())
			case col_Golem:
				rating.IceGuardian = sanitizeRating(cS.Text())
			case col_ArenaDef:
				rating.ArenaDef = sanitizeRating(cS.Text())
			case col_ArenaOff:
				rating.ArenaOff = sanitizeRating(cS.Text())
			}
		})
		if champion == nil {
			errors = append(errors, fmt.Errorf("no champion named %s", name))
			return
		}
		champion.AddRating("hellhades-tier-list", &rating, 5)
		if errSanitize := champion.Sanitize(); errSanitize != nil {
			errors = append(errors, errSanitize)
		}
		champions = append(champions, champion)
	})
	for _, champion := range champions {
		errWrite := utils.WriteToFile(fmt.Sprintf("%s/docs/champions/current/%s.json", *c.DataDirectory, champion.Slug), champion)
		if errWrite != nil {
			utils.Exit(1, errWrite)
		}
	}
	if len(errors) > 0 {
		log.Printf("errors were encountered: %v\n",  errors)
	}
}

const (
	col_Name        = 0
	col_Overall     = 1
	col_ClanBoss    = 2
	col_FactionWars = 3
	col_Spider      = 4
	col_Dragon      = 5
	col_FireKnight  = 6
	col_Golem       = 7
	col_ArenaDef    = 8
	col_ArenaOff    = 9
)

func sanitizeRating(rating string) string {
	if rating == "" {
		return rating
	}
	v, err := strconv.ParseFloat(rating, 64)
	if err != nil {
		panic(err)
	}
	switch true {
	case v >= 5.0:
		return "SS"
	case v >= 4.0:
		return "S"
	case v >= 3.0:
		return "A"
	case v >= 2.0:
		return "B"
	case v >= 1.0:
		return "C"
	}
	return "D"
}

func (c *Command) requestUrl(url string) (*goquery.Document, error) {
	req, errRequest := http.NewRequest("GET", url, nil)
	if errRequest != nil {
		return nil, errRequest
	}
	resp, errResponse := http.DefaultClient.Do(req)
	if errResponse != nil {
		return nil, errResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("request %v returned %d", req, resp.StatusCode)
	}
	return goquery.NewDocumentFromReader(resp.Body)
}
