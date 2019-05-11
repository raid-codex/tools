package website_cache_clear

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/Jeffail/tunny"
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

func (c *Command) Run() {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://raid-codex.com/?action=wpfastestcache&type=clearcacheandminified&token=%s", *c.CacheToken),
		nil,
	)
	if err != nil {
		utils.Exit(1, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Exit(1, err)
	} else if resp.StatusCode != 200 {
		utils.Exit(1, fmt.Errorf("got status code %d while trying to clear cache", resp.StatusCode))
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Exit(1, err)
	} else if string(content) == "Wrong token" {
		utils.Exit(1, fmt.Errorf("invalid cache token"))
	}
	smap, err := sitemap.Get("https://raid-codex.com/page-sitemap.xml", nil)
	if err != nil {
		utils.Exit(1, err)
	}
	pool := tunny.NewFunc(8, func(url interface{}) interface{} {
		v := url.(sitemap.URL)
		req, err := http.NewRequest("GET", v.Loc, nil)
		if err != nil {
			return err
		}
		fmt.Printf("refreshing %s\n", v.Loc)
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
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}(url)
	}
	wg.Wait()
}
