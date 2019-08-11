package fusions_page_generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/templatefuncs"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	FusionFile     *string
	TemplateFolder *string
	OutputFile     *string
	PageTemplate   *string
	DataDirectory  *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		FusionFile:     cmd.Flag("fusion-file", "Filename for the fusion").Required().String(),
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
	fusion, errFusion := c.getFusion()
	if errFusion != nil {
		utils.Exit(1, errFusion)
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
	extraData, errData := fusion.GetPageExtraData(*c.DataDirectory)
	if errData != nil {
		utils.Exit(1, errData)
	}
	buf := bytes.NewBufferString("")
	errTemplate := fusion.GetPageContent_Templates(templates, buf, extraData)
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

func (c *Command) getFusion() (*common.Fusion, error) {
	file, errFile := os.Open(*c.FusionFile)
	if errFile != nil {
		return nil, errors.Annotate(errFile, "cannot open file")
	}
	defer file.Close()

	var fusion common.Fusion
	errJSON := json.NewDecoder(file).Decode(&fusion)
	if errJSON != nil {
		return nil, errors.Annotate(errJSON, "cannot unmarshal file")
	}
	return &fusion, nil
}
