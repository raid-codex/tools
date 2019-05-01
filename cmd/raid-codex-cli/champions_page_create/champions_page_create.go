package champions_page_create

import (
	"github.com/raid-codex/tools/utils"
	"encoding/json"
	"fmt"
	"os"
	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionFile *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		ChampionFile: cmd.Flag("champion-file", "Filename for the champion").Required().String(),
	}
}

func (c *Command) Run() {
	client := utils.GetWPClient()

	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
	}
	_, errPage := utils.GetPageFromSlug(client, champion.GetPageSlug())
	if errPage != nil && !errors.IsNotFound(errPage) {
		utils.Exit(1, errPage)
	} else if errPage != nil && errors.IsNotFound(errPage) {
		errCreate := utils.CreatePage(client, champion)
		if errCreate != nil {
			utils.Exit(1, errCreate)
		}
	} else {
		fmt.Println("page already exists, ignoring")
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