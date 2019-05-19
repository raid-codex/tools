package champions_parser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

// https://spreadsheets.google.com/feeds/download/spreadsheets/Export?key=1jdrS8mnsITEWL1qREShSG3xNOZKYJuL5dUnNrUWQIjw&exportFormat=csv

type Command struct {
	CSVFile       *string
	CurrentFolder *string
	TargetFolder  *string
	NoCurrent     *bool
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		CSVFile:       cmd.Flag("csv-file", "CSV File").Required().String(),
		CurrentFolder: cmd.Flag("current-folder", "Folder where current champions are stored").String(),
		TargetFolder:  cmd.Flag("target-folder", "Folder in which the JSON files with champions data should be created").Required().String(),
		NoCurrent:     cmd.Flag("no-current", "Should be set to true if you don't want to specify any current folder, this is to avoid potential issues during champions parsing").Bool(),
	}
	return command
}

func (c *Command) Run() {
	if *c.NoCurrent == false && *c.CurrentFolder == "" {
		utils.Exit(1, errors.New("if no current folder, then --no-current should be set"))
	}
	content, errContent := c.getSourceFileContent()
	if errContent != nil {
		utils.Exit(1, errors.Annotate(errContent, "cannot read file"))
	}
	errExport := c.exportContent(content)
	if errExport != nil {
		utils.Exit(1, errors.Annotate(errExport, "cannot export content"))
	}
}

func (c *Command) debug() string {
	return fmt.Sprintf("csv-file: %s, current-folder: %s, target-folder: %s, no-current: %t", *c.CSVFile, *c.CurrentFolder, *c.TargetFolder, *c.NoCurrent)
}

type Champions struct {
	Champions common.ChampionList
}

func (c *Command) getSourceFileContent() (*Champions, error) {
	file, errFile := os.Open(*c.CSVFile)
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
		champion, err := c.getChampionByName(line[1])
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

func (c *Command) getChampionByName(name string) (*common.Champion, error) {
	nameOk, errSanitize := common.GetSanitizedName(name)
	if errSanitize != nil {
		return nil, errSanitize
	}
	if *c.CurrentFolder == "" {
		return &common.Champion{
			Name: nameOk,
		}, nil
	}
	file, errOpen := os.Open(fmt.Sprintf("%s/%s.json", *c.CurrentFolder, common.GetLinkNameFromSanitizedName(nameOk)))
	if errOpen != nil && !isNoSuchFileOrDirectory(errOpen) {
		return nil, errOpen
	} else if errOpen != nil && isNoSuchFileOrDirectory(errOpen) {
		return &common.Champion{
			Name: nameOk,
		}, nil
	}
	defer file.Close()
	var champion common.Champion
	errJSON := json.NewDecoder(file).Decode(&champion)
	if errJSON != nil {
		return nil, errJSON
	}
	return &champion, nil
}

const (
	csvSafeguardOrder = `Factions,Champion,Rarity,Element,Typ,Overall,Campaign,Arena-Off,Arena-Deff,CB (- GS),CB (+GS),IceG,Dragon,Spider,FK,Mino,Force,Magic,Spirit,Void`
)

func (c *Command) exportContent(content *Champions) error {
	content.Champions.Sort()
	for _, champion := range content.Champions {
		filename := champion.Filename()
		errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *c.TargetFolder, filename), champion)
		if errWrite != nil {
			return errWrite
		}
	}
	return nil
}
