package ps

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

type SchemaGenerator struct {
	app         *pocketbase.PocketBase
	collections []*models.Collection
}

func (api *API) GetCollections() ([]*models.Collection, error) {
	txDao := api.app.Dao()

	collections := []*models.Collection{}
	if err := txDao.CollectionQuery().OrderBy("created ASC").All(&collections); err != nil {
		return nil, err
	}

	return collections, nil
}

type Any map[string]interface{}

type Gen struct {
	T          schema.FieldOptions
	Collection *models.Collection
	File       io.Writer
}

func (g *Gen) InitOptions() string {
	requiredFields := GetRequiredFields(g.Collection.Schema)

	typeDsl := fmt.Sprintf(`

	var %s = dsl.Type("%s", func() {

		%s

		dsl.Required(%#v...)
	})
	`, Title(g.Collection.Name), g.Collection.Name, g.ParseSchema(), requiredFields)

	return typeDsl
}

func GetRequiredFields(s schema.Schema) []string {
	var requiredFields []string
	for _, field := range s.Fields() {
		if field.Required {
			requiredFields = append(requiredFields, field.Name)
		}
	}
	return requiredFields
}

func (g *Gen) ParseSchema() string {

	fields := g.Collection.Schema.Fields()

	var schemaDsl string
	for _, field := range fields {
		schemaDsl += fmt.Sprintf("\n%s\n", ParseSchemaField(field))
	}
	return schemaDsl
}

func ParseSchemaField(s *schema.SchemaField) string {
	var dslType string
	switch s.Type {
	case schema.FieldTypeText:
		dslType = "dsl.String"

	case schema.FieldTypeNumber:
		dslType = "dsl.Float64"

	case schema.FieldTypeBool:
		dslType = "dsl.Boolean"

	case schema.FieldTypeEmail:
		dslType = "dsl.String"

	case schema.FieldTypeUrl:
		dslType = "dsl.String"

	case schema.FieldTypeDate:
		// todo: add datetype
		dslType = "dsl.String"

	case schema.FieldTypeSelect:
		// todo: add options with arrayof
		dslType = "dsl.String"

	case schema.FieldTypeJson:
		dslType = "dsl.Any"

	case schema.FieldTypeFile:
		dslType = "dsl.String"

	case schema.FieldTypeRelation:
		dslType = "dsl.String"

	case schema.FieldTypeUser:
		dslType = "dsl.String"

	default:
		fmt.Println("Missing or unknown field field type.")
		return ""
	}
	return fmt.Sprintf(`dsl.Attribute("%s", %s)`, s.Name, dslType)
}

func Title(s string) string {
	//lint:ignore SA1019 update this later
	return strings.Title(s)
}

func NewFile(name string) (*Gen, error) {
	file, err := os.OpenFile(name, os.O_CREATE, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	return &Gen{
		File: file,
	}, nil
}
