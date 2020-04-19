package champions_parse_tierlist_hellhades

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

// https://docs.google.com/spreadsheets/d/1YjETkvBMVKZr7CPDjL_iIy_Wa6-psHob6so6fgKLX7c/edit#gid=0
// https://spreadsheets.google.com/feeds/download/spreadsheets/Export?key=1YjETkvBMVKZr7CPDjL_iIy_Wa6-psHob6so6fgKLX7c&exportFormat=csv
// https://docs.google.com/spreadsheets/d/1YjETkvBMVKZr7CPDjL_iIy_Wa6-psHob6so6fgKLX7c/htmlview?sle=true#gid=0
// https://spreadsheets.google.com/feeds/download/spreadsheets/Export?key=1YjETkvBMVKZr7CPDjL_iIy_Wa6-psHob6so6fgKLX7c&exportFormat=csv

type Command struct {
	CSVFile       *string
	DataDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		DataDirectory: cmd.Flag("data-directory", "Data directory").Required().String(),
		CSVFile:       cmd.Flag("csv-file", "CSV File").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	errInit := common.InitFactory(*c.DataDirectory)
	if errInit != nil {
		utils.Exit(1, errInit)
	}
	file, errFile := os.Open(*c.CSVFile)
	if errFile != nil {
		utils.Exit(1, errFile)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	content, errRead := reader.ReadAll()
	if errRead != nil {
		utils.Exit(1, errRead)
	}
	champions := make([]*common.Champion, 0)
	errs := make([]error, 0)
	for idx, line := range content {
		if idx == 0 {
			ok := false
			for _, sg := range safeGuard {
				if strings.Join(line, ",") == sg {
					ok = true
					break
				}
			}
			if !ok {
				utils.Exit(1, fmt.Errorf("invalid first line: %s", strings.Join(line, ",")))
			}
			continue
		}
		champion, errChampion := c.getChampion(line[1])
		if errChampion != nil {
			errs = append(errs, errChampion)
			continue
		}
		rating := &common.Rating{}
		rating.ArenaDef = sanitizeRating(line[rating_ArenaDef])
		rating.ArenaOff = sanitizeRating(line[rating_ArenaOff])
		rating.ClanBossWoGS = sanitizeRating(line[rating_ClanBossWOGS])
		rating.ClanBosswGS = sanitizeRating(line[rating_ClanBossWGS])
		rating.IceGuardian = sanitizeRating(line[rating_IceGolem])
		rating.Dragon = sanitizeRating(line[rating_Dragon])
		rating.Spider = sanitizeRating(line[rating_Spider])
		rating.FireKnight = sanitizeRating(line[rating_FireKnight])
		champion.AddRating("hellhades-tier-list", rating, 5)
		champions = append(champions, champion)
	}
	if len(errs) > 0 {
		utils.Exit(1, fmt.Errorf("got multiple errors: %+v", errs))
	}
	for _, champion := range champions {
		if err := champion.Sanitize(); err != nil {
			utils.Exit(1, fmt.Errorf("cannot sanitize champion %s: %s", champion.Name, err))
		}
	}
	for _, champion := range champions {
		if err := c.saveChampion(champion); err != nil {
			utils.Exit(1, fmt.Errorf("cannot save champion %s: %s", champion.Name, err))
		}
	}
}

var (
	intToRank = map[int]string{
		5: "SS",
		4: "S",
		3: "A",
		2: "B",
		1: "C",
		0: "D",
	}
	championReplacement = map[string]string{
		"Allure":               "Alure",
		"Lutheia":              "Luthiea",
		"InfernalBaroness":     "Infernal Baroness",
		"Flesh Tearer":         "Flesh-Tearer",
		"Cannoness":            "Canoness",
		"Woad Painted":         "Woad-Painted",
		"Bad-el-Kazaar":        "Bad-el-Kazar",
		"Big 'Un":              "Big'Un",
		"Teela Groremane":      "Teela Goremane",
		"Painkeeper":           "Pain Keeper",
		"Bloodhord":            "Bloodhorn",
		"Amaratine Skeleton":   "Amarantine Skeleton",
		"Lameller":             "Lamellar",
		"Steadfast Marshall":   "Steadfast Marshal",
		"Furystroker":          "Furystoker",
		"Tormenter":            "Tormentor",
		"Houndspawn":           "Hound Spawn",
		"Arablaster":           "Arbalester",
		"Siphi the lost bride": "Siphi the Lost Bride",
		"Rotos the lost groom": "Rotos the Lost Groom",
		"Ursine IceCrusher":    "Ursine Icecrusher",
		"Lordly Legionaire":    "Lordly Legionary",
	}
)

func sanitizeRating(rating string) string {
	if rating == "" {
		return rating
	}
	intV, err := strconv.Atoi(rating)
	if err != nil {
		panic(err)
	}
	return intToRank[intV]
}

func (c *Command) getChampionWithFilter(filter common.ChampionFilter) (*common.Champion, error) {
	champions, err := common.GetChampions(filter)
	if err != nil {
		return nil, err
	} else if len(champions) != 1 {
		return nil, fmt.Errorf("found %d champions", len(champions))
	}
	return champions[0], nil
}

func (c *Command) getChampion(name string) (*common.Champion, error) {
	nameOk, errSanitize := common.GetSanitizedName(name)
	if errSanitize != nil {
		return nil, errSanitize
	}
	if v, ok := championReplacement[nameOk]; ok {
		nameOk = v
	}
	champion, err := c.getChampionWithFilter(common.FilterChampionName(nameOk))
	if err != nil {
		champion, err = c.getChampionWithFilter(common.FilterChampionSlug(common.GetLinkNameFromSanitizedName(nameOk)))
		if err != nil {
			return nil, fmt.Errorf("error while looking up champion %s: %v", nameOk, err)
		}
	}
	file, errFile := os.Open(fmt.Sprintf("%s/docs/champions/current/%s.json", *c.DataDirectory, champion.Slug))
	if errFile != nil {
		return nil, errFile
	}
	defer file.Close()
	champion = &common.Champion{}
	errDecode := json.NewDecoder(file).Decode(champion)
	if errDecode != nil {
		return nil, errDecode
	}
	return champion, nil
}

func (c *Command) saveChampion(champion *common.Champion) error {
	filename := champion.Filename()
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/docs/champions/current/%s", *c.DataDirectory, filename), champion)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

var (
	safeGuard = []string{
		`Rank,Champion,Faction,Affinity,Clan Boss,Dragon,Spider,Ice Golem,FireKnight,Arena Off,Arena Def,Total Score,Average,Guide`,
		`Overall,Champion,Faction,Affinity,Clan Boss,Dragon,Spider,Ice Golem,FireKnight,Arena Off,Arena Def,Total Score,Average,Guide`,
		`Overall,Champion,Faction,Affinity,Clan Boss,Dragon,Spider,Ice Golem,FireKnight,Arena Off,Arena Def,Total Score,Average,Guide,,,,,,,,,,,`,
		`Rank,Champion,Faction,Affinity,Clan Boss,Dragon,Spider,Ice Golem,FireKnight,Arena Off,Arena Def,Total Score,Average,Guide,,,,,,,,,,,`,
	}
)

const (
	rating_ArenaOff     = 9
	rating_ArenaDef     = 10
	rating_ClanBossWOGS = 4
	rating_ClanBossWGS  = 4
	rating_IceGolem     = 7
	rating_Dragon       = 5
	rating_Spider       = 6
	rating_FireKnight   = 8
)
