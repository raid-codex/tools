package status_effects_sanitize

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	StatusEffectFile *string
	DataDirectory    *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		StatusEffectFile: cmd.Flag("status-effect-file", "Filename for the status effect").Required().String(),
		DataDirectory:    cmd.Flag("data-directory", "Data directory").Required().String(),
	}
}

func (c *Command) Run() {
	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	effect, errEffect := c.getEffect()
	if errEffect != nil {
		utils.Exit(1, errEffect)
	}
	errSanitize := effect.Sanitize()
	if errSanitize != nil {
		utils.Exit(1, errSanitize)
	}
	errWrite := utils.WriteToFile(*c.StatusEffectFile, effect)
	if errWrite != nil {
		utils.Exit(1, errWrite)
	}
}

func (c *Command) getEffect() (*common.StatusEffect, error) {
	file, errFile := os.Open(*c.StatusEffectFile)
	if errFile != nil {
		return nil, errors.Annotate(errFile, "cannot open file")
	}
	defer file.Close()

	var effect common.StatusEffect
	errJSON := json.NewDecoder(file).Decode(&effect)
	if errJSON != nil {
		return nil, errors.Annotate(errJSON, "cannot unmarshal file")
	}
	return &effect, nil
}
