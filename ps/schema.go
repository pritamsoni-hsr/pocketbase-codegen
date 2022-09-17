package ps

import (
	"encoding/json"
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

func (be *API) BuildSchema() error {
	txDao := be.app.Dao()

	existingCollections := []*models.Collection{}
	if err := txDao.CollectionQuery().OrderBy("created ASC").All(&existingCollections); err != nil {
		return err
	}

	var apiMessages []Any

	var srv string

	sc := make(Any)

	for _, _collection := range existingCollections {
		srv += ""

		sc["collectionId"] = _collection.Id
		sc["collectionName"] = _collection.Name
		sc["collectionSystem"] = _collection.System
		sc_schema := make(Any)
		sc["$schema"] = &sc_schema

		for idx, _field := range _collection.Schema.Fields() {
			_field.InitOptions()
			options, _ := json.Marshal(_field.Options)
			sc_schema[_field.Name] = Any{
				"Id":       _field.Id,
				"Name":     _field.Name,
				"Type":     _field.Type,
				"Options":  string(options),
				"Required": _field.Required,
			}
			jsonSchema, _ := json.Marshal(sc)
			apiMessages = append(apiMessages, sc)
			filename := fmt.Sprintf("message_%d.json", idx)
			if err := os.WriteFile(filename, jsonSchema, 0644); err != nil {
				return fmt.Errorf("failed to save messages %w", err)
			}
		}
	}

	fmt.Println(apiMessages)
	return nil
}

type Any map[string]interface{}

type Gen struct {
	T          schema.FieldOptions
	Collection *models.Collection
	File       io.Writer
}

func (g *Gen) InitOptions() string {
	requiredFields := GetRequiredFields(g.Collection.Schema)

	d := fmt.Sprintf(`

	var %s = dsl.Type("%s", func() {

	`, strings.Title(g.Collection.Name), g.Collection.Name)

	d += g.FieldOptions()

	d += fmt.Sprintf(`
		dsl.Required(%#v...)
	})
	`, requiredFields)

	return d
}

func (g *Gen) FieldOptions() string {

	f := g.Collection.Schema.Fields()

	var dd string
	for _, field := range f {
		dd += ParseSchemaField(field)
	}
	return dd
}

func (w *Gen) WriteLn(format string, a ...any) string {
	s := fmt.Sprintf(format, a...) + "\n"
	p := []byte(s)
	w.File.Write(p)
	return s
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

func GetRequiredFields(s schema.Schema) []string {
	var requiredFields []string
	for _, field := range s.Fields() {
		if field.Required {
			requiredFields = append(requiredFields, field.Name)
		}
	}
	return requiredFields
}

func ParseSchemaField(s *schema.SchemaField) string {
	switch s.Type {
	case schema.FieldTypeText:
		r := s.Options.(*schema.TextOptions)
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String, func() {
			dsl.MinLength(%d)
			dsl.MaxLength(%d)
			dsl.Pattern("%s")
		})
		`, s.Name, r.Min, r.Max, r.Pattern)

	case schema.FieldTypeNumber:
		r := s.Options.(*schema.NumberOptions)
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.Float64, func() {
			dsl.Minimum(%d)
			dsl.Maximum(%d)
		})`, s.Name, r.Min, r.Max)

	case schema.FieldTypeBool:
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.Boolean, func() {

		})
		`, s.Name)

	case schema.FieldTypeEmail:
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String, func() {

		})
		`, s.Name)

	case schema.FieldTypeUrl:
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String, func() {

		})
		`, s.Name)

	case schema.FieldTypeDate:
		r := s.Options.(*schema.DateOptions)
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String, func() {
			dsl.Minimum(%d)
			dsl.Maximum(%d)
		})
		`, s.Name, r.Min, r.Max)

	case schema.FieldTypeSelect:
		r := s.Options.(*schema.SelectOptions)
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String, func() {
			dsl.ArrayOf(%v)
		})
		`, s.Name, r.Values)

	case schema.FieldTypeJson:
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.Any)
		`, s.Name)

	case schema.FieldTypeFile:
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String)
		`, s.Name)

	case schema.FieldTypeRelation:
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String)
		`, s.Name)

	case schema.FieldTypeUser:
		return fmt.Sprintf(`
		dsl.Attribute("%s", dsl.String)
		`, s.Name)

	default:
		fmt.Println("Missing or unknown field field type.")
	}
	return ""
}
