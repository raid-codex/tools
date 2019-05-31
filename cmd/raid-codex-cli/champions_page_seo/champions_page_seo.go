package champions_page_seo

import (
	"encoding/json"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/seo/seowp"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionFile *string
	SetDefault   *bool
	champion     *common.Champion
	action       string
}

func New(cmd *kingpin.CmdClause, action string) *Command {
	command := &Command{
		action:       action,
		ChampionFile: cmd.Flag("champion-file", "JSON file for the Champion. Data will be edited in place if needed").Required().String(),
	}
	return command
}

func (c *Command) Run() {
	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
	}
	c.champion = champion
	switch c.action {
	case "set-default":
		c.setDefault()
	case "apply":
		c.apply()
	default:
		utils.Exit(1, errors.NotImplementedf("action %s", c.action))
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

func (c *Command) setDefault() {
	c.champion.DefaultSEO()
	errWrite := utils.WriteToFile(*c.ChampionFile, c.champion)
	if errWrite != nil {
		utils.Exit(1, errors.Annotate(errWrite, "cannot set default seo on champion"))
	}
}

func (c *Command) apply() {
	errApply := seowp.Apply(c.champion.Slug, c.champion.SEO)
	if errApply != nil {
		utils.Exit(1, errApply)
	}
}
