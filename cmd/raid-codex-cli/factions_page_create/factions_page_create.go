package factions_page_create

import (
	"github.com/raid-codex/tools/utils"
	"encoding/json"
	"fmt"
	"os"
	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	FactionFile *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		FactionFile: cmd.Flag("faction-file", "Filename for the faction").Required().String(),
	}
}

func (c *Command) Run() {
	client := utils.GetWPClient()

	faction, errFaction := c.getFaction()
	if errFaction != nil {
		utils.Exit(1, errFaction)
	}
	_, errPage := utils.GetPageFromSlug(client, faction.GetPageSlug())
	if errPage != nil && !errors.IsNotFound(errPage) {
		utils.Exit(1, errPage)
	} else if errPage != nil && errors.IsNotFound(errPage) {
		errCreate := utils.CreatePage(client, faction)
		if errCreate != nil {
			utils.Exit(1, errCreate)
		}
	} else {
		fmt.Println("page already exists, ignoring")
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