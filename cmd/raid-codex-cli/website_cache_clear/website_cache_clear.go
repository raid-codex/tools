package website_cache_clear

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{}
	return command
}

var (
	zoneID string = os.Getenv("CF_ZONE_ID")
)

func (c *Command) Run() {
	err := c.run()
	if err != nil {
		utils.Exit(1, err)
	}
}

func (c *Command) run() error {
	cf, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		return err
	}
	log.Println("purging Cloudflare cache")
	resp, err := cf.PurgeCache(context.TODO(), zoneID, cloudflare.PurgeCacheRequest{
		Everything: true,
	})
	if err != nil {
		return err
	} else if !resp.Success {
		return fmt.Errorf("could not purge cache")
	}
	log.Println("cache purged")
	return nil
}
