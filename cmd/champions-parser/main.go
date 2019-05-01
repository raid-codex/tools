package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	sourceFile    = kingpin.Arg("source-file", "CSV File with champions list and grades").Required().String()
	targetFolder  = kingpin.Arg("target-folder", "Folder in which to create the JSON files").Required().String()
	currentFolder = kingpin.Arg("current-folder", "Folder in which current champions are stored").Required().String()
)

func main() {
	kingpin.Parse()
	content, errContent := getSourceFileContent()
	if errContent != nil {
		fmt.Fprintf(os.Stderr, "cannot read file: %s\n", errContent)
		os.Exit(1)
	}
	errExport := exportContent(content)
	if errExport != nil {
		fmt.Fprintf(os.Stderr, "error while exporting: %s\n", errExport)
		os.Exit(1)
	}
}

type Champions struct {
	Champions []*common.Champion
}

func getSourceFileContent() (*Champions, error) {
	file, errFile := os.Open(*sourceFile)
	if errFile != nil {
		return nil, errFile
	}
	defer file.Close()
	champions := &Champions{}
	reader := csv.NewReader(file)
	content, errRead := reader.ReadAll()
	if errRead != nil {
		return nil, errRead
	}
	for idx, line := range content {
		if len(line) != 20 {
			return nil, fmt.Errorf("line %s has %d parts, not 20", strings.Join(line, ","), len(line))
		} else if line[0] == "Factions" {
			if strings.Join(line, ",") != csvSafeguardOrder {
				return nil, fmt.Errorf("invalid csv")
			}
			// this is a header line, skip it
			continue
		} else if line[0] == "" && idx > 0 {
			// in case faction is not mentionned on every line, use the one from previous line
			line[0] = content[idx-1][0]
		}
		champion, err := getChampionByName(line[1])
		if err != nil {
			return nil, err
		}
		champion.Faction = common.Faction{Name: line[0]}
		champion.Name = line[1]
		champion.Rarity = line[2]
		champion.Element = line[3]
		champion.Type = line[4]
		champion.Rating.Overall = line[5]
		champion.Rating.Campaign = line[6]
		champion.Rating.ArenaOff = line[7]
		champion.Rating.ArenaDef = line[8]
		champion.Rating.ClanBossWoGS = line[9]
		champion.Rating.ClanBosswGS = line[10]
		champion.Rating.IceGuardian = line[11]
		champion.Rating.Dragon = line[12]
		champion.Rating.Spider = line[13]
		champion.Rating.FireKnight = line[14]
		champion.Rating.Minotaur = line[15]
		champion.Rating.ForceDungeon = line[16]
		champion.Rating.MagicDungeon = line[17]
		champion.Rating.SpiritDungeon = line[18]
		champion.Rating.VoidDungeon = line[19]
		errSanitize := champion.Sanitize()
		if errSanitize != nil {
			return nil, errSanitize
		}
		champions.Champions = append(champions.Champions, champion)
	}
	return champions, nil
}

func isNoSuchFileOrDirectory(err error) bool {
	return strings.Contains(err.Error(), "no such file or directory")
}

func getChampionByName(name string) (*common.Champion, error) {
	nameOk, errSanitize := common.GetSanitizedName(name)
	if errSanitize != nil {
		return nil, errSanitize
	}
	file, errOpen := os.Open(fmt.Sprintf("%s/%s.json", *currentFolder, common.GetLinkNameFromSanitizedName(nameOk)))
	if errOpen != nil && !isNoSuchFileOrDirectory(errOpen) {
		return nil, errOpen
	} else if errOpen != nil && isNoSuchFileOrDirectory(errOpen) {
		return &common.Champion{
			Name: nameOk,
		}, nil
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
	csvSafeguardOrder = `Factions,Champion,Rarity,Element,Typ,Overall,Campaign,Arena-Off,Arena-Deff,CB (- GS),CB (+GS),IceG,Dragon,Spider,FK,Mino,Force,Magic,Spirit,Void`
)

func exportContent(content *Champions) error {
	for _, champion := range content.Champions {
		filename := champion.Filename()
		errWrite := writeToFile(filename, champion)
		if errWrite != nil {
			return errWrite
		}
	}
	errWrite := writeToFile("index.json", content.Champions)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func writeToFile(filename string, val interface{}) error {
	return utils.WriteToFile(fmt.Sprintf("%s/%s", *targetFolder, filename), val)
}
