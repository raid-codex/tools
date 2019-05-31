package factions_page_seo

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/seo/seowp"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	FactionFile *string
	SetDefault  *bool
	faction     *common.Faction
	action      string
}

func New(cmd *kingpin.CmdClause, action string) *Command {
	command := &Command{
		action:      action,
		FactionFile: cmd.Flag("faction-file", "JSON file for the Faction. Data will be edited in place if needed").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	faction, errFaction := c.getFaction()
	if errFaction != nil {
		utils.Exit(1, errFaction)
	}
	c.faction = faction
	switch c.action {
	case "set-default":
		c.setDefault()
	case "apply":
		c.apply()
	default:
		utils.Exit(1, errors.NotImplementedf("action %s", c.action))
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

func (c *Command) setDefault() {
	c.faction.DefaultSEO()
	errWrite := utils.WriteToFile(*c.FactionFile, c.faction)
	if errWrite != nil {
		utils.Exit(1, errors.Annotate(errWrite, "cannot set default seo on faction"))
	}
}

func (c *Command) apply() {
	errApply := seowp.Apply(c.faction.Slug, c.faction.SEO)
	if errApply != nil {
		utils.Exit(1, errApply)
	}
}
