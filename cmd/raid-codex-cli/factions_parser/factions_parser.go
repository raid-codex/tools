package factions_parser

import (
	"fmt"
	"os"
	"encoding/json"
	"sort"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
)

type Command struct {
	ChampionsDirectory *string	
	CurrentFolder *string
	TargetFolder *string
	NoCurrent *bool
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		ChampionsDirectory: cmd.Flag("champions-directory", "Folder in which current champions are stored").Required().String(),
		TargetFolder: cmd.Flag("target-folder", "Folder in which to create the JSON files").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	champions, errChampions := c.fetchChampions()
	if errChampions != nil {
		utils.Exit(1, errChampions)
	}
	factions := map[string]*common.Faction{}
	for _, champion := range champions {
		if factions[champion.Faction.Name] == nil {
			factions[champion.Faction.Name] = &champion.Faction
		}
	}

	factionsList := []*common.Faction{}
	for _, faction := range factions {
		errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *c.TargetFolder, faction.Filename()), faction)
		if errWrite != nil {
			utils.Exit(1, errWrite)
		}
		factionsList = append(factionsList, faction)
	}
	sort.SliceStable(factionsList, func (i, j int) bool {
		return factionsList[i].Slug < factionsList[j].Slug
	})
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/index.json", *c.TargetFolder), factionsList)
	if errWrite != nil {
		utils.Exit(1, errWrite)
	}
}

func (c *Command) fetchChampions() ([]*common.Champion, error) {
	file, errOpen := os.Open(fmt.Sprintf("%s/index.json", *c.ChampionsDirectory))
	if errOpen != nil {
		return nil, errOpen
	}
	var champions []*common.Champion
	errJSON := json.NewDecoder(file).Decode(&champions)
	if errJSON != nil {
		return nil, errJSON
	}
	return champions, nil
}