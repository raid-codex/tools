package fusions_sanitize

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	FusionFile    *string
	DataDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		FusionFile:    cmd.Flag("fusion-file", "Filename for the fusion").Required().String(),
		DataDirectory: cmd.Flag("data-directory", "Data directory").Required().String(),
	}
}

func (c *Command) Run() {
	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	fusion, errFusion := c.getFusion()
	if errFusion != nil {
		utils.Exit(1, errFusion)
	}
	errSanitize := fusion.Sanitize()
	if errSanitize != nil {
		utils.Exit(1, errSanitize)
	}
	errWrite := utils.WriteToFile(*c.FusionFile, fusion)
	if errWrite != nil {
		utils.Exit(1, errWrite)
	}
}

func (c *Command) getFusion() (*common.Fusion, error) {
	file, errFile := os.Open(*c.FusionFile)
	if errFile != nil {
		return nil, errors.Annotate(errFile, "cannot open file")
	}
	defer file.Close()

	var fusion common.Fusion
	errJSON := json.NewDecoder(file).Decode(&fusion)
	if errJSON != nil {
		return nil, errors.Annotate(errJSON, "cannot unmarshal file")
	}
	return &fusion, nil
}
