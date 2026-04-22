APP_NAME := hxtrainer
MAIN_PKG := ./cmd/hxtrainer

.PHONY: test build install fmt release-snapshot

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) $(MAIN_PKG)

install:
	go install $(MAIN_PKG)

fmt:
	gofmt -w ./cmd ./internal ./challenges_data

release-snapshot:
	goreleaser release --snapshot --clean
