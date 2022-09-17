format:
	find ./api -type f  -exec sed  -i '' 's/dsl\.//g'  {} +
	gofmt -w ./api
