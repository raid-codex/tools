package champions_schema_validate

import (
	"fmt"

	"github.com/raid-codex/tools/utils"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	ChampionFile *string
	SchemaFile   *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		ChampionFile: cmd.Flag("champion-file", "Filename for the champion").Required().String(),
		SchemaFile:   cmd.Flag("schema-file", "Filename for the schema").Required().String(),
	}
}

func (c *Command) Run() {
	schemaLoader := gojsonschema.NewReferenceLoader(*c.SchemaFile)
	documentLoader := gojsonschema.NewReferenceLoader(*c.ChampionFile)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		utils.Exit(1, err)
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
