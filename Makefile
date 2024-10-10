.PHONY: all
all: fmt test build

.PHONY: fmt
fmt:
	@gofmt -w ./**/*.go

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -v ./database/repo

.PHONY: clean
clean:
	@go clean

.PHONY: run
run:
	@go run main.go

.PHONY: replicate
replicate:
	@litestream replicate -config litestream.yml

.PHONY: restore
restore:
	@litestream restore -if-db-not-exists -config litestream.yml ${DB_PATH}

.PHONY: exec
exec:
	@litestream replicate -config litestream.yml -exec "$(CURDIR)/referrals"
