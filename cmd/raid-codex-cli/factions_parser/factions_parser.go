package factions_parser

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionsDirectory *string
	TargetFolder       *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		ChampionsDirectory: cmd.Flag("champions-directory", "Folder in which current champions are stored").Required().String(),
		TargetFolder:       cmd.Flag("target-folder", "Folder in which to create the JSON files").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	champions, errChampions := c.fetchChampions()
	if errChampions != nil {
		utils.Exit(1, errChampions)
	}
	championsPerFaction := map[string]int64{}
	factions := map[string]*common.Faction{}
	for _, champion := range champions {
		if factions[champion.Faction.Name] == nil {
			factions[champion.Faction.Name] = &champion.Faction
		}
		championsPerFaction[champion.Faction.Name] += 1
	}

	factionsList := []*common.Faction{}
	for _, faction := range factions {
		faction.NumberOfChampions = championsPerFaction[faction.Name]
		errSanitize := faction.Sanitize()
		if errSanitize != nil {
			utils.Exit(1, errSanitize)
		}
		errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *c.TargetFolder, faction.Filename()), faction)
		if errWrite != nil {
			utils.Exit(1, errWrite)
		}
		factionsList = append(factionsList, faction)
	}
	sort.SliceStable(factionsList, func(i, j int) bool {
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
	defer file.Close()
	var champions []*common.Champion
	errJSON := json.NewDecoder(file).Decode(&champions)
	if errJSON != nil {
		return nil, errJSON
	}
	return champions, nil
}
