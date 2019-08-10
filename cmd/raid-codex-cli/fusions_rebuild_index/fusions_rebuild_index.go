package fusions_rebuild_index

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
	FusionsDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		FusionsDirectory: cmd.Flag("fusion-directory", "Folder in which fusions are stored").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	fusions, errChampions := c.fetchFusions()
	if errChampions != nil {
		utils.Exit(1, errChampions)
	}
	fusions.Sort()
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/index.json", *c.FusionsDirectory), fusions)
	if errWrite != nil {
		utils.Exit(1, errors.Annotate(errWrite, "cannot write to file"))
	}
}

func (c *Command) fetchFusions() (common.FusionList, error) {
	dir, err := ioutil.ReadDir(*c.FusionsDirectory)
	if err != nil {
		return nil, errors.Annotatef(err, "cannot open %s", *c.FusionsDirectory)
	}
	var statusEffects common.FusionList
	for _, file := range dir {
		if file.Name() == "index.json" {
			continue
		} else if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		var fusion common.Fusion
		err := func() error {
			filename := fmt.Sprintf("%s/%s", *c.FusionsDirectory, file.Name())
			f, errOpen := os.Open(filename)
			if errOpen != nil {
				return errors.Annotatef(errOpen, "cannot open %s", filename)
			}
			defer f.Close()
			errJSON := json.NewDecoder(f).Decode(&fusion)
			return errJSON
		}()
		if err != nil {
			return nil, err
		}
		statusEffects = append(statusEffects, &fusion)
	}
	return statusEffects, nil
}
