In progress

## Swagger openapi schema generator for pocketbase

- [x] generate all messages from `db/_collections/*/schema`
	- [ ] add formats like datetime, file, json, jsonschema
- [ ] get all api endpoints available to only users, we will would admin types for now.
- [ ] build openapi spec from the above definitions using go-swagger/ or other swagger gen or goa.design/goa

---

Each schema will be versioned, and schema is only available for a running app.

History of schema changes will not be available.
