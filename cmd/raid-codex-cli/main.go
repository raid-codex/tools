package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_characteristics_parser"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_page_generate"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_page_seo"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_parse_tierlist"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_parse_tierlist_hellhades"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_parser"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_rate"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_rebuild_index"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_sanitize"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/champions_video_add"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_page_generate"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_page_seo"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_parser"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_rebuild_index"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/factions_sanitize"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/fusions_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/fusions_page_generate"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/fusions_rebuild_index"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/fusions_sanitize"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/parse_full_sheet"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/schema_validate"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/scrap_ayumilove_champions"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/scrap_gameronion_champions"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/scrap_raidshadowlegendspro_champions"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/scrap_wikia_characteristics"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/server_run"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/status_effects_page_create"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/status_effects_page_generate"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/status_effects_rebuild_index"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/status_effects_sanitize"
	"github.com/raid-codex/tools/cmd/raid-codex-cli/website_cache_clear"
	"github.com/raid-codex/tools/utils"
	_ "github.com/raid-codex/tools/utils/logger" // init logger
	"gopkg.in/alecthomas/kingpin.v2"
)

type Runnable interface {
	Run()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "recovered from panic while running command '%s'\n%v\n%s\n", strings.Join(os.Args, " "), r, string(debug.Stack()))
		}
	}()
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

	championsVideo       = champions.Command("video", "Video")
	championsVideoAdd    = championsVideo.Command("add", "Add video to champion")
	championsVideoAddCmd = champions_video_add.New(championsVideoAdd)

	championsRating                 = champions.Command("rating", "Rate champion")
	championsRatingAddFromSource    = championsRating.Command("add-from-source", "Add rating from source")
	championsRatingAddFromSourceCmd = champions_rate.New(championsRatingAddFromSource)

	championsSanitize    = champions.Command("sanitize", "sanitize champion file")
	championsSanitizeCmd = champions_sanitize.New(championsSanitize)

	championsParser    = champions.Command("parser", "parse champions from csv file")
	championsParserCmd = champions_parser.New(championsParser)

	championsCharacteristics          = champions.Command("characteristics", "deal with champions characteristics")
	championsCharacteristicsParser    = championsCharacteristics.Command("parse", "parse champions characteristics")
	championsCharacteristicsParserCmd = champions_characteristics_parser.New(championsCharacteristicsParser)

	championsRebuildIndex    = champions.Command("rebuild-index", "rebuild champions index")
	championsRebuildIndexCmd = champions_rebuild_index.New(championsRebuildIndex)

	championsPage          = champions.Command("page", "Handle champion page")
	championsPageCreate    = championsPage.Command("create", "Create the page for the champion")
	championsPageCreateCmd = champions_page_create.New(championsPageCreate)

	championsPageGenerate    = championsPage.Command("generate", "Generate HTML for Champion page")
	championsPageGenerateCmd = champions_page_generate.New(championsPageGenerate)

	championsPageSeo              = championsPage.Command("seo", "Deal with SEO for a champion page")
	championsPageSeoSetDefault    = championsPageSeo.Command("set-default", "Reset SEO settings to default")
	championsPageSeoSetDefaultCmd = champions_page_seo.New(championsPageSeoSetDefault, "set-default")
	championsPageSeoApply         = championsPageSeo.Command("apply", "Apply SEO settings to champion page")
	championsPageSeoApplyCmd      = champions_page_seo.New(championsPageSeoApply, "apply")

	championsParse            = champions.Command("parse", "Parse stuff about champions")
	championsParseTierList    = championsParse.Command("tier-list", "Tier list")
	championsParseTierListCmd = champions_parse_tierlist.New(championsParseTierList)

	championsParseTierListHellhades    = championsParse.Command("tier-list-hellhades", "Hellhades tier list")
	championsParseTierListHellhadesCmd = champions_parse_tierlist_hellhades.New(championsParseTierListHellhades)

	factions = app.Command("factions", "do stuff with factions")

	factionsSanitize    = factions.Command("sanitize", "sanitize faction file")
	factionsSanitizeCmd = factions_sanitize.New(factionsSanitize)

	factionsParse    = factions.Command("parse", "parse factions from champions json files")
	factionsParseCmd = factions_parser.New(factionsParse)

	factionsRebuildIndex    = factions.Command("rebuild-index", "rebuild faction index")
	factionsRebuildIndexCmd = factions_rebuild_index.New(factionsRebuildIndex)

	factionsPage          = factions.Command("page", "Handle faction page")
	factionsPageCreate    = factionsPage.Command("create", "Create the page for the faction")
	factionsPageCreateCmd = factions_page_create.New(factionsPageCreate)

	factionsPageGenerate    = factionsPage.Command("generate", "Generate the page for the faction")
	factionsPageGenerateCmd = factions_page_generate.New(factionsPageGenerate)

	factionsPageSeo              = factionsPage.Command("seo", "Deal with SEO for a faction page")
	factionsPageSeoSetDefault    = factionsPageSeo.Command("set-default", "Reset SEO settings to default")
	factionsPageSeoSetDefaultCmd = factions_page_seo.New(factionsPageSeoSetDefault, "set-default")
	factionsPageSeoApply         = factionsPageSeo.Command("apply", "Apply SEO settings to faction page")
	factionsPageSeoApplyCmd      = factions_page_seo.New(factionsPageSeoApply, "apply")

	scrap = app.Command("scrap", "Scrap stuff from the internet")

	scrapWikiaCharacteristics    = scrap.Command("wikia-characteristics", "Scrap data from wikia characteristics")
	scrapWikiaCharacteristicsCmd = scrap_wikia_characteristics.New(scrapWikiaCharacteristics)

	scrapGameronion = scrap.Command("gameronion", "Scrap data from gameronion")

	scrapGameronionChampions    = scrapGameronion.Command("champions", "Scrap champions")
	scrapGameronionChampionsCmd = scrap_gameronion_champions.New(scrapGameronionChampions)

	scrapAyumilove = scrap.Command("ayumilove", "Scrap data from ayumilove")

	scrapAyumiloveChampions    = scrapAyumilove.Command("champions", "Scrap champions")
	scrapAyumiloveChampionsCmd = scrap_ayumilove_champions.New(scrapAyumiloveChampions)

	scrapRaidShadowLegendsPro = scrap.Command("raidshadowlegendspro", "Scrap data from raidshadowlegends.pro")

	scrapRaidShadowLegendsProChampions    = scrapRaidShadowLegendsPro.Command("champions", "Scrap champions")
	scrapRaidShadowLegendsProChampionsCmd = scrap_raidshadowlegendspro_champions.New(scrapRaidShadowLegendsProChampions)

	website              = app.Command("website", "Stuff for website")
	websiteCache         = website.Command("cache", "Stuff with website cache")
	websiteCacheClear    = websiteCache.Command("clear", "Clear cache of website")
	websiteCacheClearCmd = website_cache_clear.New(websiteCacheClear)

	statusEffect = app.Command("status-effect", "Stuff for status effect")

	statusEffectSanitize    = statusEffect.Command("sanitize", "Sanitize a status effect file")
	statusEffectSanitizeCmd = status_effects_sanitize.New(statusEffectSanitize)

	statusEffectRebuildIndex    = statusEffect.Command("rebuild-index", "Rebuild status effects index")
	statusEffectRebuildIndexCmd = status_effects_rebuild_index.New(statusEffectRebuildIndex)

	statusEffectPage            = statusEffect.Command("page", "Handle status effect page")
	statusEffectPageGenerate    = statusEffectPage.Command("generate", "Generate HTML for status effect page")
	statusEffectPageGenerateCmd = status_effects_page_generate.New(statusEffectPageGenerate)

	statusEffectPageCreate    = statusEffectPage.Command("create", "Create or update page on the website")
	statusEffectPageCreateCmd = status_effects_page_create.New(statusEffectPageCreate)

	schema            = app.Command("schema", "Stuff for schemas")
	schemaValidate    = schema.Command("validate", "Validate a file against a schema")
	schemaValidateCmd = schema_validate.New(schemaValidate)

	parse             = app.Command("parse", "Parse stuff")
	parseFullSheet    = parse.Command("full-sheet", "Parse the full-sheet stuff")
	parseFullSheetCmd = parse_full_sheet.New(parseFullSheet)

	fusions = app.Command("fusions", "do stuff with fusions")

	fusionsSanitize    = fusions.Command("sanitize", "sanitize fusion file")
	fusionsSanitizeCmd = fusions_sanitize.New(fusionsSanitize)

	fusionsRebuildIndex    = fusions.Command("rebuild-index", "rebuild faction index")
	fusionsRebuildIndexCmd = fusions_rebuild_index.New(fusionsRebuildIndex)

	fusionsPage            = fusions.Command("page", "do stuff with fusion pages")
	fusionsPageGenerate    = fusionsPage.Command("generate", "generate page for fusion")
	fusionsPageGenerateCmd = fusions_page_generate.New(fusionsPageGenerate)

	fusionsPageCreate    = fusionsPage.Command("create", "create page for fusion")
	fusionsPageCreateCmd = fusions_page_create.New(fusionsPageCreate)

	server = app.Command("server", "Server")

	serverRun    = server.Command("run", "Run the server")
	serverRunCmd = server_run.New(serverRun)

	runByCmd = map[string]Runnable{
		"champions rating add-from-source":     championsRatingAddFromSourceCmd,
		"champions parser":                     championsParserCmd,
		"champions parse tier-list":            championsParseTierListCmd,
		"champions parse tier-list-hellhades":  championsParseTierListHellhadesCmd,
		"factions parse":                       factionsParseCmd,
		"champions page create":                championsPageCreateCmd,
		"factions page create":                 factionsPageCreateCmd,
		"scrap wikia-characteristics":          scrapWikiaCharacteristicsCmd,
		"scrap gameronion champions":           scrapGameronionChampionsCmd,
		"scrap ayumilove champions":            scrapAyumiloveChampionsCmd,
		"scrap raidshadowlegendspro champions": scrapRaidShadowLegendsProChampionsCmd,
		"champions page seo set-default":       championsPageSeoSetDefaultCmd,
		"champions page seo apply":             championsPageSeoApplyCmd,
		"champions rebuild-index":              championsRebuildIndexCmd,
		"factions page seo set-default":        factionsPageSeoSetDefaultCmd,
		"factions page seo apply":              factionsPageSeoApplyCmd,
		"factions page generate":               factionsPageGenerateCmd,
		"champions characteristics parse":      championsCharacteristicsParserCmd,
		"champions sanitize":                   championsSanitizeCmd,
		"champions video add":                  championsVideoAddCmd,
		"factions sanitize":                    factionsSanitizeCmd,
		"website cache clear":                  websiteCacheClearCmd,
		"champions page generate":              championsPageGenerateCmd,
		"schema validate":                      schemaValidateCmd,
		"status-effect sanitize":               statusEffectSanitizeCmd,
		"status-effect rebuild-index":          statusEffectRebuildIndexCmd,
		"factions rebuild-index":               factionsRebuildIndexCmd,
		"status-effect page generate":          statusEffectPageGenerateCmd,
		"status-effect page create":            statusEffectPageCreateCmd,
		"parse full-sheet":                     parseFullSheetCmd,
		"fusions sanitize":                     fusionsSanitizeCmd,
		"fusions rebuild-index":                fusionsRebuildIndexCmd,
		"fusions page generate":                fusionsPageGenerateCmd,
		"fusions page create":                  fusionsPageCreateCmd,
		"server run":                           serverRunCmd,
	}
)
