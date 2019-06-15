package status_effects_page_create

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"github.com/raid-codex/tools/utils/wp"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	StatusEffectFile *string
	TemplateFile     *string
	DataDirectory    *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		StatusEffectFile: cmd.Flag("status-effect-file", "Filename for the status effect").Required().String(),
		TemplateFile:     cmd.Flag("template-file", "Template file").Required().String(),
		DataDirectory:    cmd.Flag("data-directory", "Data directory").Required().String(),
	}
}

func (c *Command) Run() {
	client := wp.GetWPClient()

	effect, errEffect := c.getStatusEffect()
	if errEffect != nil {
		utils.Exit(1, errEffect)
	} else if strings.HasSuffix(effect.Slug, "-2") {
		utils.Exit(0, fmt.Errorf("skipping, since it's an upgrade of %s", effect.Slug[:len(effect.Slug)-2]))
	}
	page, errPage := wp.GetPageFromSlug(client, effect.GetPageSlug())
	if errPage != nil && !errors.IsNotFound(errPage) {
		utils.Exit(1, errPage)
	} else if errPage != nil && errors.IsNotFound(errPage) {
		errCreate := wp.CreatePage(client, effect, *c.TemplateFile, *c.DataDirectory, nil)
		if errCreate != nil {
			utils.Exit(1, errCreate)
		}
	} else {
		errUpdate := wp.UpdatePage(client, page, effect, *c.TemplateFile, *c.DataDirectory, nil)
		if errUpdate != nil {
			utils.Exit(1, errUpdate)
		}
	}
}

func (c *Command) getStatusEffect() (*common.StatusEffect, error) {
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
