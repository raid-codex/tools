package champions_rate

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	DataDirectory *string
	ChampionName  *string
	Source        *string
	Weight        *int
	Campaign      *string
	ClanBoss      *string
	FireKnight    *string
	Arena         *string
	Dragon        *string
	IceGolem      *string
	Spider        *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		DataDirectory: cmd.Flag("data-directory", "Data directory").Required().String(),
		ChampionName:  cmd.Flag("champion-name", "Champion name").Required().String(),
		Source:        cmd.Flag("source", "Source").Required().String(),
		Weight:        cmd.Flag("weight", "Weight").Required().Int(),
		Campaign:      cmd.Flag("campaign", "").String(),
		ClanBoss:      cmd.Flag("clan-boss", "").String(),
		FireKnight:    cmd.Flag("fire-knight", "").String(),
		Arena:         cmd.Flag("arena", "").String(),
		Dragon:        cmd.Flag("dragon", "").String(),
		IceGolem:      cmd.Flag("ice-golem", "").String(),
		Spider:        cmd.Flag("spider", "").String(),
	}
	return command
}

func (c *Command) Run() {
	if errInit := common.InitFactory(*c.DataDirectory); errInit != nil {
		utils.Exit(1, errInit)
	}
	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
	}
	var rating *common.Rating
	for _, r := range champion.AllRatings {
		if r.Source == *c.Source {
			rating = r.Rating
			break
		}
	}
	if rating == nil {
		rating = &common.Rating{}
	}
	if *c.Campaign != "" {
		rating.Campaign = *c.Campaign
	}
	if *c.ClanBoss != "" {
		rating.ClanBosswGS = *c.ClanBoss
		rating.ClanBossWoGS = *c.ClanBoss
	}
	if *c.Arena != "" {
		rating.ArenaOff = *c.Arena
		rating.ArenaDef = *c.Arena
	}
	if *c.IceGolem != "" {
		rating.IceGuardian = *c.IceGolem
	}
	if *c.FireKnight != "" {
		rating.FireKnight = *c.FireKnight
	}
	if *c.Spider != "" {
		rating.Spider = *c.Spider
	}
	if *c.Dragon != "" {
		rating.Dragon = *c.Dragon
	}

	champion.AddRating(*c.Source, rating, *c.Weight)
	if errSanitize := champion.Sanitize(); errSanitize != nil {
		utils.Exit(1, errSanitize)
	}

	if errSave := c.saveChampion(champion); errSave != nil {
		utils.Exit(1, errSave)
	}
}

func (c *Command) getChampion() (*common.Champion, error) {
	nameOk, errSanitize := common.GetSanitizedName(*c.ChampionName)
	if errSanitize != nil {
		return nil, errSanitize
	}
	champions, err := common.GetChampions(common.FilterChampionName(nameOk))
	if err != nil {
		return nil, err
	} else if len(champions) != 1 {
		return nil, fmt.Errorf("found %d champions with name %s", len(champions), nameOk)
	}
	file, errFile := os.Open(fmt.Sprintf("%s/docs/champions/current/%s.json", *c.DataDirectory, champions[0].Slug))
	if errFile != nil {
		return nil, errFile
	}
	defer file.Close()
	var champion common.Champion
	errDecode := json.NewDecoder(file).Decode(&champion)
	if errDecode != nil {
		return nil, errDecode
	}
	return &champion, nil
}

func (c *Command) saveChampion(champion *common.Champion) error {
	filename := champion.Filename()
	errWrite := utils.WriteToFile(fmt.Sprintf("%s/docs/champions/current/%s", *c.DataDirectory, filename), champion)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
