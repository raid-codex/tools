package factions_page_create

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/templatefuncs"
	"github.com/raid-codex/tools/utils"
	"github.com/raid-codex/tools/utils/wp"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	FactionFile    *string
	DataDirectory  *string
	TemplateFolder *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		FactionFile:    cmd.Flag("faction-file", "Filename for the faction").Required().String(),
		TemplateFolder: cmd.Flag("template-folder", "Template folder").Required().String(),
		DataDirectory:  cmd.Flag("data-directory", "Data directory").Required().String(),
	}
}

func (c *Command) Run() {
	client := wp.GetWPClient()

	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	tmpl, errTmpl := c.loadTemplates()
	if errTmpl != nil {
		utils.Exit(1, errTmpl)
	}
	faction, errFaction := c.getFaction()
	if errFaction != nil {
		utils.Exit(1, errFaction)
	}
	page, errPage := wp.GetPageFromSlug(client, faction.GetPageSlug())
	if errPage != nil && !errors.IsNotFound(errPage) {
		utils.Exit(1, errPage)
	} else if errPage != nil && errors.IsNotFound(errPage) {
		errCreate := wp.CreatePage(client, faction, "", *c.DataDirectory, tmpl)
		if errCreate != nil {
			utils.Exit(1, errCreate)
		}
	} else {
		errUpdate := wp.UpdatePage(client, page, faction, "", *c.DataDirectory, tmpl)
		if errUpdate != nil {
			utils.Exit(1, errUpdate)
		}
	}
}

func (c *Command) getFaction() (*common.Faction, error) {
	file, errFile := os.Open(*c.FactionFile)
	if errFile != nil {
		return nil, errors.Annotate(errFile, "cannot open file")
	}
	defer file.Close()

	var faction common.Faction
	errJSON := json.NewDecoder(file).Decode(&faction)
	if errJSON != nil {
		return nil, errors.Annotate(errJSON, "cannot unmarshal file")
	}
	return &faction, nil
}

func (c *Command) loadTemplates() (*template.Template, error) {
	files, errFiles := ioutil.ReadDir(*c.TemplateFolder)
	if errFiles != nil {
		return nil, errFiles
	}
	templateFiles := make([]string, 0)
	for _, file := range files {
		templateFiles = append(templateFiles, fmt.Sprintf("%s/%s", *c.TemplateFolder, file.Name()))
	}
	return template.New("main.html").Funcs(templatefuncs.FuncMap).ParseFiles(templateFiles...)
}
