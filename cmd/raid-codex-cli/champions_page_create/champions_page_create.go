package champions_page_create

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
	ChampionFile   *string
	TemplateFolder *string
	DataDirectory  *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		ChampionFile:   cmd.Flag("champion-file", "Filename for the champion").Required().String(),
		TemplateFolder: cmd.Flag("template-folder", "Template folder").String(),
		DataDirectory:  cmd.Flag("data-directory", "Data directory").String(),
	}
}

func (c *Command) Run() {
	client := wp.GetWPClient()
	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
	}
	content := fmt.Sprintf(`[raid-codex-champion-page slug="%s"]`, champion.GetPageSlug())
	page, errPage := wp.GetPageFromSlug(client, champion.GetPageSlug())
	if errPage != nil && !errors.IsNotFound(errPage) {
		utils.Exit(1, errPage)
	} else if errPage != nil && errors.IsNotFound(errPage) {
		errCreate := wp.CreatePage_Content(client, champion, content)
		if errCreate != nil {
			utils.Exit(1, errCreate)
		}
	} else {
		errUpdate := wp.UpdatePage_Content(client, page, champion, content)
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
