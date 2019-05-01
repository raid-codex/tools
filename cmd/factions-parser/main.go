package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	championsDirectory = kingpin.Arg("champions-directory", "Folder in which current champions are stored").Required().String()
	targetFolder       = kingpin.Arg("target-folder", "Folder in which to create the JSON files").Required().String()
)

func main() {
	kingpin.Parse()
	champions, errChampions := fetchChampions()
	if errChampions != nil {
		exit(errChampions)
	}
	factions := map[string]*common.Faction{}
	for _, champion := range champions {
		if factions[champion.Faction.Name] == nil {
			factions[champion.Faction.Name] = &champion.Faction
		}
	}

	factionsList := []*common.Faction{}
	for _, faction := range factions {
		errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *targetFolder, faction.Filename()), faction)
		if errWrite != nil {
			exit(errWrite)
		}
		factionsList = append(factionsList, faction)
	}
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/index.json", *targetFolder), factionsList)
	if errWrite != nil {
		exit(errWrite)
	}
}

func fetchChampions() ([]*common.Champion, error) {
	file, errOpen := os.Open(fmt.Sprintf("%s/index.json", *championsDirectory))
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

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}
