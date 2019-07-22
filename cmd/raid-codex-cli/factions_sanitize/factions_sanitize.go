package factions_sanitize

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	FactionFile   *string
	DataDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		FactionFile:   cmd.Flag("faction-file", "Filename for the faction").Required().String(),
		DataDirectory: cmd.Flag("data-directory", "Directory containing all the game data").Required().String(),
	}
}

func (c *Command) Run() {
	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	faction, errFaction := c.getFaction()
	if errFaction != nil {
		utils.Exit(1, errFaction)
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
