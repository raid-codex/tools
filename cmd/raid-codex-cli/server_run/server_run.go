package server_run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/templatefuncs"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	DataDirectory  *string
	TemplateFolder *string
	PageTemplate   *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		DataDirectory:  cmd.Flag("data-directory", "Directory containing data").Required().String(),
		TemplateFolder: cmd.Flag("template-folder", "Template folder").Required().String(),
		PageTemplate:   cmd.Flag("page-template", "Page template file").Required().String(),
	}
}

func (c *Command) Run() {
	errFactory := common.InitFactory(*c.DataDirectory)
	if errFactory != nil {
		utils.Exit(1, errFactory)
	}
	srv := gin.New()
	srv.Use(errorHandler)
	srv.GET("/web/champions/:champion_slug", c.webChampionSlug)
	if err := srv.Run(":8080"); err != nil {
		utils.Exit(1, err)
	}
}

func errorHandler(ctx *gin.Context) {
	ctx.Next()
	detectedErrors := ctx.Errors.String()
	log.Printf("errors: %s\n", detectedErrors)
}

func (c *Command) webChampionSlug(ctx *gin.Context) {
	slug := ctx.Param("champion_slug")
	champion, errChampion := c.getChampion(slug)
	if errChampion != nil {
		ctx.AbortWithError(500, errChampion)
		return
	}
	templates, errLoad := c.loadTemplates("champion")
	if errLoad != nil {
		ctx.AbortWithError(500, errLoad)
		return
	}
	extraData, errData := champion.GetPageExtraData(*c.DataDirectory)
	if errData != nil {
		ctx.AbortWithError(500, errData)
		return
	}
	buf := bytes.NewBufferString("")
	errTemplate := champion.GetPageContent_Templates(templates, buf, extraData)
	if errTemplate != nil {
		ctx.AbortWithError(500, errTemplate)
		return
	}
	pageTemplate, errPageTemplate := ioutil.ReadFile(*c.PageTemplate)
	if errPageTemplate != nil {
		ctx.AbortWithError(500, errPageTemplate)
		return
	}
	buf2 := bytes.NewBufferString("")
	tmpl, errTmpl := template.New("page").Funcs(templatefuncs.FuncMap).Parse(string(pageTemplate))
	if errTmpl != nil {
		ctx.AbortWithError(500, errTmpl)
		return
	}
	errExecute := tmpl.Execute(buf2, map[string]interface{}{"Page": buf.String()})
	if errExecute != nil {
		ctx.AbortWithError(500, errExecute)
		return
	}
	ctx.Data(200, "text/html", buf2.Bytes())
}

func (c *Command) loadTemplates(dir string) (*template.Template, error) {
	dir = fmt.Sprintf("%s/%s", *c.TemplateFolder, dir)
	files, errFiles := ioutil.ReadDir(dir)
	if errFiles != nil {
		return nil, errFiles
	}
	templateFiles := make([]string, 0)
	for _, file := range files {
		templateFiles = append(templateFiles, fmt.Sprintf("%s/%s", dir, file.Name()))
	}
	return template.New("main.html").Funcs(templatefuncs.FuncMap).ParseFiles(templateFiles...)
}

func (c *Command) getChampion(slug string) (*common.Champion, error) {
	file, errFile := os.Open(fmt.Sprintf("%s/docs/champions/current/%s.json", *c.DataDirectory, slug))
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
