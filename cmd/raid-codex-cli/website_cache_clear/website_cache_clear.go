package website_cache_clear

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Jeffail/tunny"
	"github.com/cloudflare/cloudflare-go"
	"github.com/raid-codex/tools/utils"
	"github.com/yterajima/go-sitemap"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	CacheToken *string
}

func New(cmd *kingpin.CmdClause) *Command {
	command := &Command{
		CacheToken: cmd.Flag("cache-token", "Token used to refresh cache").Required().String(),
	}
	return command
}

var (
	zoneID string = os.Getenv("CF_ZONE_ID")
	ruleID string = os.Getenv("CF_PAGE_RULE_ID")
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
	pageRule, err := cf.PageRule(zoneID, ruleID)
	if err != nil {
		return err
	}
	pageRule.Status = "disabled"
	log.Println("disabling page rule")
	err = cf.UpdatePageRule(zoneID, ruleID, pageRule)
	if err != nil {
		return err
	}
	defer func() {
		log.Println("re-enabling page rule")
		pageRule.Status = "active"
		err = cf.UpdatePageRule(zoneID, ruleID, pageRule)
		if err != nil {
			log.Printf("error while reactivating page rule: %v\n", err)
		}
	}()
	log.Println("purging Cloudflare cache")
	resp, err := cf.PurgeCache(zoneID, cloudflare.PurgeCacheRequest{
		Everything: true,
	})
	if err != nil {
		return err
	} else if !resp.Success {
		return fmt.Errorf("could not purge cache")
	}
	err = c.emptyWPcache()
	if err != nil {
		return err
	}
	err = c.refreshWP()
	if err != nil {
		return err
	}
	return nil
}

func (c *Command) refreshWP() error {
	log.Println("refreshing website")
	smap, err := sitemap.Get("https://raid-codex.com/page-sitemap.xml", nil)
	if err != nil {
		return err
	}
	pool := tunny.NewFunc(8, func(url interface{}) interface{} {
		v := url.(sitemap.URL)
		req, err := http.NewRequest("GET", v.Loc, nil)
		if err != nil {
			return err
		}
		log.Printf("refreshing %s\n", v.Loc)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		} else if resp.StatusCode != 200 {
			return fmt.Errorf("url %s returned status code %d", v.Loc, resp.StatusCode)
		}
		return nil
	})
	wg := sync.WaitGroup{}
	for _, url := range smap.URL {
		wg.Add(1)
		go func(u sitemap.URL) {
			defer wg.Done()
			err := pool.Process(u)
			if err != nil {
				log.Printf("error while refreshing %s: %v\n", u, err)
			}
		}(url)
	}
	wg.Wait()
	log.Println("website refreshed")
	return nil
}

func (c *Command) emptyWPcache() error {
	log.Println("Requesting WP cache to be emptied")
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://raid-codex.com/?action=wpfastestcache&type=clearcacheandminified&token=%s", *c.CacheToken),
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("got status code %d while trying to clear cache", resp.StatusCode)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	} else if string(content) == "Wrong token" {
		return fmt.Errorf("invalid cache token")
	}
	log.Println("WP cache emptied")
	return nil
}
