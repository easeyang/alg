.PHONY: build
build:
	echo "Building..."
	GOOS=linux go build -o bin/$(app) cmd/$(app)/main.go
