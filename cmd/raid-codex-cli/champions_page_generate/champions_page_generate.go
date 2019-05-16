package champions_page_generate

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
	ChampionFile *string
	TemplateFile *string
	OutputFile   *string
	PageTemplate *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		ChampionFile: cmd.Flag("champion-file", "Filename for the champion").Required().String(),
		TemplateFile: cmd.Flag("template-file", "Template file").Required().String(),
		OutputFile:   cmd.Flag("output-file", "Output file").Required().String(),
		PageTemplate: cmd.Flag("page-template", "Page template file").Required().String(),
	}
}

func (c *Command) Run() {
	champion, errChampion := c.getChampion()
	if errChampion != nil {
		utils.Exit(1, errChampion)
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
	buf := bytes.NewBufferString("")
	errTemplate := champion.GetPageContent(inputFile, buf)
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
