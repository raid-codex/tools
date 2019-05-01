package main

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

// URL: https://raid-shadow-legends.fandom.com/wiki/Tier_List ; GET CSV FROM : http://www.convertcsv.com/html-table-to-csv.htm

var (
	sourceFile    = kingpin.Arg("source-file", "CSV File with champions list and characteristics").Required().String()
	currentFolder = kingpin.Arg("current-folder", "Folder in which current champions are stored. JSON files will be edited in-place").Required().String()
)

func main() {
	kingpin.Parse()
	file, errFile := os.Open(*sourceFile)
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
		if len(line) != 18 {
			utils.Exit(1, fmt.Errorf("line %s has %d parts, not 18", strings.Join(line, ","), len(line)))
		}
		champion, err := getChampionByName(line[1])
		if err != nil {
			utils.Exit(1, err)
		}
		characteristics := champion.Characteristics[60]
		characteristics.HP = mustInt64(line[5])
		characteristics.Attack = mustInt64(line[6])
		characteristics.Defense = mustInt64(line[7])
		characteristics.Speed = mustInt64(line[8])
		characteristics.CriticalRate = float64(mustInt64(line[9][0:len(line[9])-1])) / 100.0
		characteristics.CriticalDamage = float64(mustInt64(line[10][0:len(line[10])-1])) / 100.0
		characteristics.Resistance = mustInt64(line[11])
		characteristics.Accuracy = mustInt64(line[12])
		champion.Characteristics[60] = characteristics
		errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *currentFolder, champion.Filename()), champion)
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

func getChampionByName(name string) (*common.Champion, error) {
	nameOk, errSanitize := common.GetSanitizedName(name)
	if errSanitize != nil {
		return nil, errSanitize
	}
	file, errOpen := os.Open(fmt.Sprintf("%s/%s.json", *currentFolder, common.GetLinkNameFromSanitizedName(nameOk)))
	if errOpen != nil {
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
	csvSafeguardOrder = `Rarity,Element,Tribe,Role,Name,HP,ATK,DEF,SPD,C-Rate,C-Dmg,RES,ACC,Skill1,Skill2,Skill3,Skill4,Aura`
)
