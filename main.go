package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type API struct {
	app     *pocketbase.PocketBase
	version string
}

func main() {
	app := pocketbase.New()
	api := &API{
		app:     app,
		version: "0.0.1",
	}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/openapi",
			Handler: api.GetSchema,
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/openapi/validate",
			Handler: api.ValidateSchema,
		})
		return nil
	})
	app.OnBeforeServe().Add(func(data *core.ServeEvent) error {
		return api.BuildSchema()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func (be *API) BuildSchema() error {
	txDao := be.app.Dao()

	existingCollections := []*models.Collection{}
	if err := txDao.CollectionQuery().OrderBy("created ASC").All(&existingCollections); err != nil {
		return err
	}

	var apiMessages []Any

	sc := make(Any)

	for _, _collection := range existingCollections {
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
			fmt.Println(string(jsonSchema))
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

func (be *API) GetSchema(c echo.Context) error {
	return c.String(200, "WIP")
}

func (be *API) ValidateSchema(c echo.Context) error {
	value := c.QueryParam("version")
	if value == "" {
		return c.String(400, "correct select version")
	}
	if be.version == c.FormValue("version") {
		return c.String(200, "you are using the latest version")
	}
	return c.String(200, "WIP")
}

type Any map[string]interface{}
