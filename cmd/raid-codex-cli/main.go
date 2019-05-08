package main

import (
	"os"

	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_characteristics_parser"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_page_seo"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_parser"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_rebuild_index"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_sanitize"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_schema_validate"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_page_seo"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_parser"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/scrap_wikia_characteristics"
	"github.com/raid-codex/tools/utils"
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

	championsSanitize    = champions.Command("sanitize", "sanitize champion file")
	championsSanitizeCmd = champions_sanitize.New(championsSanitize)

	championsParse    = champions.Command("parse", "parse champions from csv file")
	championsParseCmd = champions_parser.New(championsParse)

	championsCharacteristics          = champions.Command("characteristics", "deal with champions characteristics")
	championsCharacteristicsParser    = championsCharacteristics.Command("parse", "parse champions characteristics")
	championsCharacteristicsParserCmd = champions_characteristics_parser.New(championsCharacteristicsParser)

	championsRebuildIndex    = champions.Command("rebuild-index", "rebuild champions index")
	championsRebuildIndexCmd = champions_rebuild_index.New(championsRebuildIndex)

	championsPage          = champions.Command("page", "Handle champion page")
	championsPageCreate    = championsPage.Command("create", "Create the page for the champion")
	championsPageCreateCmd = champions_page_create.New(championsPageCreate)

	championsPageSeo              = championsPage.Command("seo", "Deal with SEO for a champion page")
	championsPageSeoSetDefault    = championsPageSeo.Command("set-default", "Reset SEO settings to default")
	championsPageSeoSetDefaultCmd = champions_page_seo.New(championsPageSeoSetDefault, "set-default")
	championsPageSeoApply         = championsPageSeo.Command("apply", "Apply SEO settings to champion page")
	championsPageSeoApplyCmd      = champions_page_seo.New(championsPageSeoApply, "apply")

	championsSchema            = champions.Command("schema", "Handle champion schema")
	championsSchemaValidate    = championsSchema.Command("validate", "Validate a champion against its schema")
	championsSchemaValidateCmd = champions_schema_validate.New(championsSchemaValidate)

	factions = app.Command("factions", "do stuff with factions")

	factionsParse    = factions.Command("parse", "parse factions from champions json files")
	factionsParseCmd = factions_parser.New(factionsParse)

	factionsPage          = factions.Command("page", "Handle faction page")
	factionsPageCreate    = factionsPage.Command("create", "Create the page for the faction")
	factionsPageCreateCmd = factions_page_create.New(factionsPageCreate)

	factionsPageSeo              = factionsPage.Command("seo", "Deal with SEO for a faction page")
	factionsPageSeoSetDefault    = factionsPageSeo.Command("set-default", "Reset SEO settings to default")
	factionsPageSeoSetDefaultCmd = factions_page_seo.New(factionsPageSeoSetDefault, "set-default")
	factionsPageSeoApply         = factionsPageSeo.Command("apply", "Apply SEO settings to faction page")
	factionsPageSeoApplyCmd      = factions_page_seo.New(factionsPageSeoApply, "apply")

	scrap = app.Command("scrap", "Scrap stuff from the internet")

	scrapWikiaCharacteristics    = scrap.Command("wikia-characteristics", "Scrap data from wikia characteristics")
	scrapWikiaCharacteristicsCmd = scrap_wikia_characteristics.New(scrapWikiaCharacteristics)

	runByCmd = map[string]Runnable{
		"champions parse":                 championsParseCmd,
		"factions parse":                  factionsParseCmd,
		"champions page create":           championsPageCreateCmd,
		"factions page create":            factionsPageCreateCmd,
		"scrap wikia-characteristics":     scrapWikiaCharacteristicsCmd,
		"champions page seo set-default":  championsPageSeoSetDefaultCmd,
		"champions page seo apply":        championsPageSeoApplyCmd,
		"champions rebuild-index":         championsRebuildIndexCmd,
		"factions page seo set-default":   factionsPageSeoSetDefaultCmd,
		"factions page seo apply":         factionsPageSeoApplyCmd,
		"champions characteristics parse": championsCharacteristicsParserCmd,
		"champions sanitize":              championsSanitizeCmd,
		"champions schema validate":       championsSchemaValidateCmd,
	}
)
