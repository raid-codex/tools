package status_effects_rebuild_index

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
	StatusEffectsDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		StatusEffectsDirectory: cmd.Flag("status-effects-directory", "Folder in which status effects are stored").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	effects, errChampions := c.fetchEffects()
	if errChampions != nil {
		utils.Exit(1, errChampions)
	}
	effects.Sort()
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/index.json", *c.StatusEffectsDirectory), effects)
	if errWrite != nil {
		utils.Exit(1, errors.Annotate(errWrite, "cannot write to file"))
	}
}

func (c *Command) fetchEffects() (common.StatusEffectList, error) {
	dir, err := ioutil.ReadDir(*c.StatusEffectsDirectory)
	if err != nil {
		return nil, err
	}
	var statusEffects common.StatusEffectList
	for _, file := range dir {
		if file.Name() == "index.json" {
			continue
		} else if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		var effect common.StatusEffect
		err := func() error {
			f, errOpen := os.Open(fmt.Sprintf("%s/%s", *c.StatusEffectsDirectory, file.Name()))
			if errOpen != nil {
				return errOpen
			}
			defer f.Close()
			errJSON := json.NewDecoder(f).Decode(&effect)
			return errJSON
		}()
		if err != nil {
			return nil, err
		}
		statusEffects = append(statusEffects, &effect)
	}
	return statusEffects, nil
}
