package factions_rebuild_index

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
	FactionsDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		FactionsDirectory: cmd.Flag("factions-directory", "Folder in which current factions are stored").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	factions, errFactions := c.fetchFactions()
	if errFactions != nil {
		utils.Exit(1, errFactions)
	}
	factions.Sort()
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/index.json", *c.FactionsDirectory), factions)
	if errWrite != nil {
		utils.Exit(1, errors.Annotate(errWrite, "cannot write to file"))
	}
}

func (c *Command) fetchFactions() (common.FactionList, error) {
	dir, err := ioutil.ReadDir(*c.FactionsDirectory)
	if err != nil {
		return nil, err
	}
	var factions common.FactionList
	for _, file := range dir {
		if file.Name() == "index.json" {
			continue
		} else if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		var faction common.Faction
		err := func() error {
			f, errOpen := os.Open(fmt.Sprintf("%s/%s", *c.FactionsDirectory, file.Name()))
			if errOpen != nil {
				return errOpen
			}
			defer f.Close()
			errJSON := json.NewDecoder(f).Decode(&faction)
			return errJSON
		}()
		if err != nil {
			return nil, err
		}
		factions = append(factions, &faction)
	}
	return factions, nil
}
