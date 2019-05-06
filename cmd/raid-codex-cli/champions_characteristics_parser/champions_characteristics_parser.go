package champions_characteristics_parser

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

// URL: https://docs.google.com/spreadsheets/d/1DvC4_OisDZXiMi2rI9nBJC8J7Oqu8n2lstvnXHCSXgg/edit#gid=0

type Command struct {
	CSVFile         *string
	ChampionsFolder *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		CSVFile:         cmd.Flag("csv-file", "CSV File with champions list and characteristics").Required().String(),
		ChampionsFolder: cmd.Flag("champions-folder", "Folder in which current champions are stored. JSON files will be edited in-place").Required().String(),
	}
}

func (c *Command) Run() {
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
	for idx, line := range content {
		if idx == 0 {
			if strings.Join(line, ",") != csvSafeguardOrder {
				utils.Exit(1, fmt.Errorf("invalid csv"))
			}
			// skip first line, it's the header
			continue
		}
		if len(line) != 17 {
			utils.Exit(1, fmt.Errorf("line %s has %d parts, not 17", strings.Join(line, ","), len(line)))
		}
		champion, err := c.getChampionByName(line[1])
		if err != nil {
			utils.Exit(1, err)
		}
		if champion.Rarity == "" {
			// it's a new champion let's do something buddy
			champion.Faction.Name = line[0]
			champion.Name = line[1]
			champion.Type = line[2]
			champion.Rarity = line[3]
			errSanitize := champion.Sanitize()
			if errSanitize != nil {
				utils.Exit(1, errSanitize)
			}
		}
		characteristics := champion.Characteristics[60]
		characteristics.HP = mustInt64(line[4])
		characteristics.Attack = mustInt64(line[5])
		characteristics.Defense = mustInt64(line[6])
		characteristics.Speed = mustInt64(line[7])
		characteristics.CriticalRate = float64(mustInt64(line[8])) / 100.0
		characteristics.CriticalDamage = float64(mustInt64(line[9])) / 100.0
		characteristics.Resistance = mustInt64(line[10])
		characteristics.Accuracy = mustInt64(line[11])
		champion.Characteristics[60] = characteristics
		errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *c.ChampionsFolder, champion.Filename()), champion)
		if errWrite != nil {
			utils.Exit(1, errWrite)
		}
	}
}

func mustInt64(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		utils.Exit(1, err)
	}
	return v
}

func (c *Command) getChampionByName(name string) (*common.Champion, error) {
	nameOk, errSanitize := common.GetSanitizedName(name)
	if errSanitize != nil {
		return nil, errSanitize
	}
	file, errOpen := os.Open(fmt.Sprintf("%s/%s.json", *c.ChampionsFolder, common.GetLinkNameFromSanitizedName(nameOk)))
	if errOpen != nil && strings.Contains(errOpen.Error(), "no such file or directory") {
		// create champion
		champion := &common.Champion{
			Slug: nameOk,
			Name: name,
		}
		return champion, nil
	} else if errOpen != nil {
		return nil, errOpen
	} else {
		defer file.Close()
		var champion common.Champion
		errJSON := json.NewDecoder(file).Decode(&champion)
		if errJSON != nil {
			return nil, errJSON
		}
		return &champion, nil
	}
}

const (
	csvSafeguardOrder = `faction,title,type,rarity,health,attack,defense,speed,crit_rate,crit_damage,resist,accuracy,skill_1,skill_2,skill_3,skill_4,skill_5`
)
