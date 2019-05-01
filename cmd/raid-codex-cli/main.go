package main

import (
	"github.com/raid-codex/tools/cmd/raid-codex-cli/scrap_wikia_characteristics"
	"os"

	"github.com/raid-codex/tools/utils"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_parser"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_parser"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Runnable interface {
	Run()
}

func main() {
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		utils.Exit(1, err)
	}
	if runnable, ok := runByCmd[cmd]; ok {
		runnable.Run()
	} else {
		app.Usage([]string{})
	}
}

var (
	app = kingpin.New("raid-codex-cli", "help")

	champions = app.Command("champions", "do stuff with champions")

	championsParse = champions.Command("parse", "parse champions from csv file")
	championsParseCmd = champions_parser.New(championsParse)

	championsPage = champions.Command("page", "Handle champion page")
	championsPageCreate = championsPage.Command("create", "Create the page for the champion")
	championsPageCreateCmd = champions_page_create.New(championsPageCreate)

	factions = app.Command("factions", "do stuff with factions")

	factionsParse = factions.Command("parse", "parse factions from champions json files")
	factionsParseCmd = factions_parser.New(factionsParse)

	factionsPage = factions.Command("page", "Handle faction page")
	factionsPageCreate = factionsPage.Command("create", "Create the page for the faction")
	factionsPageCreateCmd = factions_page_create.New(factionsPageCreate)

	scrap = app.Command("scrap", "Scrap stuff from the internet")

	scrapWikiaCharacteristics = scrap.Command("wikia-characteristics", "Scrap data from wikia characteristics")
	scrapWikiaCharacteristicsCmd = scrap_wikia_characteristics.New(scrapWikiaCharacteristics)

	runByCmd = map[string]Runnable{
		"champions parse": championsParseCmd,
		"factions parse": factionsParseCmd,
		"champions page create": championsPageCreateCmd,
		"factions page create": factionsPageCreateCmd,
		"scrap wikia-characteristics": scrapWikiaCharacteristicsCmd,
	}
)