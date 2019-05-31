package status_effects_page_generate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	StatusEffectFile *string
	TemplateFile     *string
	OutputFile       *string
	PageTemplate     *string
	DataDirectory    *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		StatusEffectFile: cmd.Flag("status-effect-file", "Filename for the status effect").Required().String(),
		DataDirectory:    cmd.Flag("data-directory", "Data directory").Required().String(),
		TemplateFile:     cmd.Flag("template-file", "Template file").Required().String(),
		OutputFile:       cmd.Flag("output-file", "Output file").Required().String(),
		PageTemplate:     cmd.Flag("page-template", "Page template file").Required().String(),
	}
}

func (c *Command) Run() {
	effect, errEffect := c.getEffect()
	if errEffect != nil {
		utils.Exit(1, errEffect)
	}
	outputFile, errOutput := os.Create(*c.OutputFile)
	if errOutput != nil {
		utils.Exit(1, errOutput)
	}
	defer outputFile.Close()
	inputFile, errInput := os.Open(*c.TemplateFile)
	if errInput != nil {
		utils.Exit(1, errInput)
	}
	defer inputFile.Close()
	extraData, errData := effect.GetPageExtraData(*c.DataDirectory)
	if errData != nil {
		utils.Exit(1, errData)
	}
	buf := bytes.NewBufferString("")
	errTemplate := effect.GetPageContent(inputFile, buf, extraData)
	if errTemplate != nil {
		utils.Exit(1, errTemplate)
	}
	pageTemplate, errPageTemplate := ioutil.ReadFile(*c.PageTemplate)
	if errPageTemplate != nil {
		utils.Exit(1, errPageTemplate)
	}
	tmpl, errTmpl := template.New("page").Parse(string(pageTemplate))
	if errTmpl != nil {
		utils.Exit(1, errTmpl)
	}
	errExecute := tmpl.Execute(outputFile, map[string]interface{}{"Page": buf.String()})
	if errExecute != nil {
		utils.Exit(1, errExecute)
	}
}

func (c *Command) getEffect() (*common.StatusEffect, error) {
	file, errFile := os.Open(*c.StatusEffectFile)
	if errFile != nil {
		return nil, errors.Annotate(errFile, "cannot open file")
	}
	defer file.Close()

	var effect common.StatusEffect
	errJSON := json.NewDecoder(file).Decode(&effect)
	if errJSON != nil {
		return nil, errors.Annotate(errJSON, "cannot unmarshal file")
	}
	return &effect, nil
}
