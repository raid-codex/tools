package champions_rebuild_index

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
	ChampionsDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		ChampionsDirectory: cmd.Flag("champions-directory", "Folder in which current champions are stored").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	champions, errChampions := c.fetchChampions()
	if errChampions != nil {
		utils.Exit(1, errChampions)
	}
	champions.Sort()
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/index.json", *c.ChampionsDirectory), champions)
	if errWrite != nil {
		utils.Exit(1, errors.Annotate(errWrite, "cannot write to file"))
	}
}

func (c *Command) fetchChampions() (common.ChampionList, error) {
	dir, err := ioutil.ReadDir(*c.ChampionsDirectory)
	if err != nil {
		return nil, err
	}
	var champions common.ChampionList
	for _, file := range dir {
		if file.Name() == "index.json" {
			continue
		} else if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		var champion common.Champion
		err := func() error {
			f, errOpen := os.Open(fmt.Sprintf("%s/%s", *c.ChampionsDirectory, file.Name()))
			if errOpen != nil {
				return errOpen
			}
			defer f.Close()
			errJSON := json.NewDecoder(f).Decode(&champion)
			return errJSON
		}()
		if err != nil {
			return nil, err
		}
		champions = append(champions, &champion)
	}
	return champions, nil
}
