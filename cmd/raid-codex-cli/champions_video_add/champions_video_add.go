package champions_video_add

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	DataDirectory *string
	ChampionSlug  *string
	Author        *string
	VideoID       *string
	Source        *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		DataDirectory: cmd.Flag("data-directory", "Data directory").Required().String(),
		ChampionSlug:  cmd.Flag("champion-slug", "Champion slug").Required().String(),
		Author:        cmd.Flag("author", "Video author").Required().String(),
		VideoID:       cmd.Flag("video-id", "Video ID").Required().String(),
		Source:        cmd.Flag("source", "Source").Required().String(),
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
	champion.Videos = append(champion.Videos, &common.Video{
		Author:    *c.Author,
		Source:    *c.Source,
		ID:        *c.VideoID,
		DateAdded: time.Now().Format(time.RFC3339),
	})
	if errSanitize := champion.Sanitize(); errSanitize != nil {
		utils.Exit(1, errSanitize)
	}
	if errSave := c.saveChampion(champion); errSave != nil {
		utils.Exit(1, errSave)
	}
}

func (c *Command) getChampion() (*common.Champion, error) {
	champions, err := common.GetChampions(common.FilterChampionSlug(*c.ChampionSlug))
	if err != nil {
		return nil, err
	} else if len(champions) != 1 {
		return nil, fmt.Errorf("found %d champions with slug %s", len(champions), *c.ChampionSlug)
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
