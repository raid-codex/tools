package champions_page_create

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionFile *string
	TemplateFile *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		ChampionFile: cmd.Flag("champion-file", "Filename for the champion").Required().String(),
		TemplateFile: cmd.Flag("template-file", "Template file").Required().String(),
	}
}

func (c *Command) Run() {
	client := utils.GetWPClient()

	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
	}
	page, errPage := utils.GetPageFromSlug(client, champion.GetPageSlug())
	if errPage != nil && !errors.IsNotFound(errPage) {
		utils.Exit(1, errPage)
	} else if errPage != nil && errors.IsNotFound(errPage) {
		errCreate := utils.CreatePage(client, champion, *c.TemplateFile)
		if errCreate != nil {
			utils.Exit(1, errCreate)
		}
	} else {
		errUpdate := utils.UpdatePage(client, page, champion, *c.TemplateFile)
		if errUpdate != nil {
			utils.Exit(1, errUpdate)
		}
	}
}

func (c *Command) getChampion() (*common.Champion, error) {
	file, errFile := os.Open(*c.ChampionFile)
	if errFile != nil {
		return nil, errors.Annotate(errFile, "cannot open file")
	}
	defer file.Close()

	var champion common.Champion
	errJSON := json.NewDecoder(file).Decode(&champion)
	if errJSON != nil {
		return nil, errors.Annotate(errJSON, "cannot unmarshal file")
	}
	return &champion, nil
}
