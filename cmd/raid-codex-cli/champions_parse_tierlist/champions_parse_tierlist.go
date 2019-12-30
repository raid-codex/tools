package champions_parse_tierlist

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

// https://spreadsheets.google.com/feeds/download/spreadsheets/Export?key=1jdrS8mnsITEWL1qREShSG3xNOZKYJuL5dUnNrUWQIjw&exportFormat=csv

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
	for idx, line := range content {
		if idx == 0 {
			if strings.Join(line, ",") != safeGuard {
				utils.Exit(1, fmt.Errorf("invalid first line: %s", strings.Join(line, ",")))
			}
			continue
		}
		champion, errChampion := c.getChampion(line[0])
		if errChampion != nil {
			utils.Exit(1, errChampion)
		}
		champion.Rating.Overall = sanitizeOverall(line[rating_Overall])
		champion.Rating.Campaign = sanitizeRating(line[rating_Campaign])
		champion.Rating.ArenaDef = sanitizeRating(line[rating_ArenaDef])
		champion.Rating.ArenaOff = sanitizeRating(line[rating_ArenaOff])
		champion.Rating.ClanBossWoGS = sanitizeRating(line[rating_ClanBossWOGS])
		champion.Rating.ClanBosswGS = sanitizeRating(line[rating_ClanBossWGS])
		champion.Rating.IceGuardian = sanitizeRating(line[rating_IceGolem])
		champion.Rating.Dragon = sanitizeRating(line[rating_Dragon])
		champion.Rating.Spider = sanitizeRating(line[rating_Spider])
		champion.Rating.FireKnight = sanitizeRating(line[rating_FireKnight])
		champion.Rating.Minotaur = sanitizeRating(line[rating_Minotaur])
		champion.Rating.ForceDungeon = sanitizeRating(line[rating_Force])
		champion.Rating.MagicDungeon = sanitizeRating(line[rating_Magic])
		champion.Rating.SpiritDungeon = sanitizeRating(line[rating_Spirit])
		champion.Rating.VoidDungeon = sanitizeRating(line[rating_Void])
		champion.Rating.FactionWars = sanitizeRating(line[rating_FactionWars])
		champions = append(champions, champion)
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

func sanitizeRating(rating string) string {
	switch rating {
	case "G":
		return "SS"
	case "T":
		return "S"
	}
	return rating
}

func sanitizeOverall(v string) string {
	intV, err := strconv.Atoi(v)
	if err != nil {
		panic(err) // fixme
	}
	switch true {
	case intV < 12:
		return "D"
	case intV < 22:
		return "C"
	case intV < 32:
		return "B"
	case intV < 45:
		return "A"
	case intV < 63:
		return "S"
	}
	return "SS"
}

func (c *Command) getChampion(name string) (*common.Champion, error) {
	nameOk, errSanitize := common.GetSanitizedName(name)
	if errSanitize != nil {
		return nil, errSanitize
	}
	champions, err := common.GetChampions(common.FilterChampionName(nameOk))
	if err != nil {
		return nil, err
	} else if len(champions) != 1 {
		return nil, fmt.Errorf("found %d champions with name %s", len(champions), nameOk)
	}
	file, errFile := os.Open(fmt.Sprintf("%s/docs/champions/current/%s.json", *c.DataDirectory, champions[0].Slug))
	if errFile != nil {
		return nil, errFile
	}
	defer file.Close()
	var champion common.Champion
	errDecode := json.NewDecoder(file).Decode(&champion)
	if errDecode != nil {
		return nil, errDecode
	}
	return &champion, nil
}

func (c *Command) saveChampion(champion *common.Champion) error {
	filename := champion.Filename()
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/docs/champions/current/%s", *c.DataDirectory, filename), champion)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

const (
	safeGuard           = `Champion,Factions,Rarity,Element,Type,Overall,Campaign,Offence,Defence,CB (-T6),CB (+T6),IceGolem,Dragon,Spider,FireKnight,Mino,Force,Magic,Spirit,Void,Factions`
	rating_Overall      = 5
	rating_Campaign     = 6
	rating_ArenaOff     = 7
	rating_ArenaDef     = 8
	rating_ClanBossWOGS = 9
	rating_ClanBossWGS  = 10
	rating_IceGolem     = 11
	rating_Dragon       = 12
	rating_Spider       = 13
	rating_FireKnight   = 14
	rating_Minotaur     = 15
	rating_Force        = 16
	rating_Magic        = 17
	rating_Spirit       = 18
	rating_Void         = 19
	rating_FactionWars  = 20
)
