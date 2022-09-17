format:
	find ./api -type f  -exec sed  -i '' 's/dsl\.//g'  {} +
	gofmt -w ./api

serve:
	go run cmd/gen.go serve


gen_openapi_spec:
	goa gen github.com/pritamsoni-hsr/pocketbase-codegen/api
	cp gen/http/openapi.json openapi.json
	rm -rf gen
	$(MAKE) _gen_typescript_client

_gen_typescript_client:
	docker run --rm \
    -v ${PWD}:/app openapitools/openapi-generator-cli generate \
    -i /app/openapi.json \
    -g typescript-fetch \
    --skip-validate-spec \
    -o /app/api/openapi
