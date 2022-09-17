package api

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("api", func() {
	Title("PocketBase API")
	Version("0.0.1")
	Server("http", func() {
		Host("development", func() {
			URI("http://localhost:8090")
		})
	})
})
