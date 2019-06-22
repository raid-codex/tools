package champions_page_generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/raid-codex/tools/templatefuncs"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionFile   *string
	TemplateFolder *string
	OutputFile     *string
	PageTemplate   *string
	DataDirectory  *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		ChampionFile:   cmd.Flag("champion-file", "Filename for the champion").Required().String(),
		DataDirectory:  cmd.Flag("data-directory", "Data directory").Required().String(),
		TemplateFolder: cmd.Flag("template-folder", "Template folder").Required().String(),
		OutputFile:     cmd.Flag("output-file", "Output file").Required().String(),
		PageTemplate:   cmd.Flag("page-template", "Page template file").Required().String(),
	}
}

func (c *Command) Run() {
	errInit := common.InitFactory(*c.DataDirectory)
	if errInit != nil {
		utils.Exit(1, errInit)
	}
	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
	}
	outputFile, errOutput := os.Create(*c.OutputFile)
	if errOutput != nil {
		utils.Exit(1, errOutput)
	}
	defer outputFile.Close()
	templates, errLoad := c.loadTemplates()
	if errLoad != nil {
		utils.Exit(1, errLoad)
	}
	extraData, errData := champion.GetPageExtraData(*c.DataDirectory)
	if errData != nil {
		utils.Exit(1, errData)
	}
	buf := bytes.NewBufferString("")
	errTemplate := champion.GetPageContent_Templates(templates, buf, extraData)
	if errTemplate != nil {
		utils.Exit(1, errTemplate)
	}
	pageTemplate, errPageTemplate := ioutil.ReadFile(*c.PageTemplate)
	if errPageTemplate != nil {
		utils.Exit(1, errPageTemplate)
	}
	tmpl, errTmpl := template.New("page").Funcs(templatefuncs.FuncMap).Parse(string(pageTemplate))
	if errTmpl != nil {
		utils.Exit(1, errTmpl)
	}
	errExecute := tmpl.Execute(outputFile, map[string]interface{}{"Page": buf.String()})
	if errExecute != nil {
		utils.Exit(1, errExecute)
	}
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
