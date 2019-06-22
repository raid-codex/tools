package champions_sanitize

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionFile  *string
	DataDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		ChampionFile:  cmd.Flag("champion-file", "Filename for the champion").Required().String(),
		DataDirectory: cmd.Flag("data-directory", "Data directory").Required().String(),
	}
}

func (c *Command) Run() {
	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
	}
	errSanitize := champion.Sanitize()
	if errSanitize != nil {
		utils.Exit(1, errSanitize)
	}
	errWrite := utils.WriteToFile(*c.ChampionFile, champion)
	if errWrite != nil {
		utils.Exit(1, errWrite)
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
