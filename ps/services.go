package ps

import (
	"bytes"
	"text/template"
)

type DslService struct {
	Name string
	Type string
}

func GenService(name string) string {

	dslService := DslService{Name: name, Type: Title(name)}

	tmpl, err := template.New("test").Parse(`

	var _ = Service("{{.Name}}", func() {

		Method("list", func() {
			Result(ArrayOf({{.Type}}))
			HTTP(func() {
				GET("/api/collections/{{.Name}}/records")
			})
		})

		Method("view", func() {
			Result({{.Type}})
			HTTP(func() {
				GET("/api/collections/{{.Name}}/records/:id")
			})
		})

		Method("create", func() {
			Payload({{.Type}})
			Result({{.Type}})
			HTTP(func() {
				POST("/api/collections/{{.Name}}/records")
			})
		})

		Method("update", func() {
			Payload({{.Type}})
			Result({{.Type}})
			HTTP(func() {
				PATCH("/api/collections/{{.Name}}/records/:id")
			})
		})

		Method("delete", func() {
			HTTP(func() {
				DELETE("/api/collections/{{.Name}}/records/:id")
			})
		})

		Files("./openapi", "{{.Name}}.json")
	})

	`)
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, dslService)
	if err != nil {
		panic(err)
	}
	return tpl.String()
}
