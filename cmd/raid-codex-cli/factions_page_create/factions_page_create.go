package factions_page_create

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"github.com/raid-codex/tools/utils/wp"
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
	client := wp.GetWPClient()

	faction, errFaction := c.getFaction()
	if errFaction != nil {
		utils.Exit(1, errFaction)
	}
	_, errPage := wp.GetPageFromSlug(client, faction.GetPageSlug())
	if errPage != nil && !errors.IsNotFound(errPage) {
		utils.Exit(1, errPage)
	} else if errPage != nil && errors.IsNotFound(errPage) {
		errCreate := wp.CreatePage(client, faction, "", "", nil)
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
