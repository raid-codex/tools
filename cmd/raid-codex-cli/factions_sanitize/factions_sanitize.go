package factions_sanitize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	FactionFile        *string
	ChampionsDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		FactionFile:        cmd.Flag("faction-file", "Filename for the faction").Required().String(),
		ChampionsDirectory: cmd.Flag("champions-directory", "Directory containing all the champions").Required().String(),
	}
}

func (c *Command) Run() {
	faction, errFaction := c.getFaction()
	if errFaction != nil {
		utils.Exit(1, errFaction)
	}
	championList, errChampions := c.fetchChampions()
	if errChampions != nil {
		utils.Exit(1, errChampions)
	}
	faction.NumberOfChampions = 0
	for _, champion := range championList {
		if champion.FactionSlug == faction.Slug {
			faction.NumberOfChampions++
		}
	}
	errSanitize := faction.Sanitize()
	if errSanitize != nil {
		utils.Exit(1, errSanitize)
	}
	errWrite := utils.WriteToFile(*c.FactionFile, faction)
	if errWrite != nil {
		utils.Exit(1, errWrite)
	}
}

func (c *Command) getFaction() (*common.Faction, error) {
	file, errFile := os.Open(*c.FactionFile)
	if errFile != nil {
		return nil, errors.Annotate(errFile, "cannot open file")
	}
	defer file.Close()

	var faction common.Faction
	errJSON := json.NewDecoder(file).Decode(&faction)
	if errJSON != nil {
		return nil, errors.Annotate(errJSON, "cannot unmarshal file")
	}
	return &faction, nil
}

func (c *Command) fetchChampions() (common.ChampionList, error) {
	dir, err := ioutil.ReadDir(*c.ChampionsDirectory)
	if err != nil {
		return nil, err
	}
	var champions common.ChampionList
	for _, file := range dir {
		if file.Name() == "index.json" {
			continue
		} else if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		var champion common.Champion
		err := func() error {
			f, errOpen := os.Open(fmt.Sprintf("%s/%s", *c.ChampionsDirectory, file.Name()))
			if errOpen != nil {
				return errOpen
			}
			defer f.Close()
			errJSON := json.NewDecoder(f).Decode(&champion)
			return errJSON
		}()
		if err != nil {
			return nil, err
		}
		champions = append(champions, &champion)
	}
	return champions, nil
}
