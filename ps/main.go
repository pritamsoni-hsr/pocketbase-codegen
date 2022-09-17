package ps

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v5"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const ApiSpecOutput = "./api/spec.go"
const OpenApiSpec = "openapi.json"

type API struct {
	app     *pocketbase.PocketBase
	version string
}

func Run() {
	api := &API{
		version: "0.0.1",
		app: pocketbase.NewWithConfig(pocketbase.Config{
			DefaultDebug:   true,
			DefaultDataDir: "pb_data.bak",
		}),
	}

	api.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		if err := api.registerRoutes(e); err != nil {
			return err
		}
		api.GenSchema()
		return nil
	})

	if err := api.app.Start(); err != nil {
		log.Fatal(err)
	}
}

func (api *API) registerRoutes(e *core.ServeEvent) error {
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
}

func (be *API) GetSchema(c echo.Context) error {
	f, err := ioutil.ReadFile(OpenApiSpec)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.String(200, string(f))
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

func (api *API) GenSchema() {
	d := Gen{}

	collections, err := api.GetCollections()
	if err != nil {
		fmt.Println(err.Error())
	}

	g := SchemaGenerator{
		app:         api.app,
		collections: collections,
	}

	tmpl := `
	package api

	import (
		. "goa.design/goa/v3/dsl"
	)

	`
	for _, col := range g.collections {
		d.Collection = col
		tmpl += d.InitOptions()
		tmpl += "\n"
		tmpl += GenService(d.Collection.Name)

	}

	ioutil.WriteFile(ApiSpecOutput, []byte(tmpl), 0644)
	cmd := exec.Command("make", "format")
	_, err = cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
