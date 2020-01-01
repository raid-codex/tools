package schema_validate

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/raid-codex/tools/utils"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	File        *string
	SchemaFile  *string
	Definitions *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		File:        cmd.Flag("file", "Filename to check").Required().String(),
		SchemaFile:  cmd.Flag("schema-file", "Filename for the schema").Required().String(),
		Definitions: cmd.Flag("definitions", "Filename for schema definitions").Required().String(),
	}
}

func (c *Command) Run() {
	definitions, err := ioutil.ReadFile(*c.Definitions)
	if err != nil {
		utils.Exit(1, err)
	}
	schema, err := ioutil.ReadFile(*c.SchemaFile)
	if err != nil {
		utils.Exit(1, err)
	}
	schemaString := strings.Replace(string(schema), `"definitions": {}`, fmt.Sprintf(`"definitions": %s`, string(definitions)), 1)
	schemaLoader := gojsonschema.NewStringLoader(schemaString)
	documentLoader := gojsonschema.NewReferenceLoader(*c.File)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		utils.Exit(1, fmt.Errorf("cannot validate: %s", err))
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		errMessage := "The document is not valid. see errors :\n"
		for _, desc := range result.Errors() {
			errMessage = fmt.Sprintf("%s - %s\n", errMessage, desc)
		}
		utils.Exit(1, fmt.Errorf(errMessage))
	}
}
