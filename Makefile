TAGS ?= ""
GO_BIN ?= "go"
GIT_REV := $(shell git describe --tags HEAD)
GO_BUILD_FLAGS := -ldflags "-X github.com/markbates/pkger.Version=dev-$(GIT_REV)"

install: tidy
	cd ./cmd/pkger && $(GO_BIN) install $(GO_BUILD_FLAGS) -tags ${TAGS} -v .
	make tidy

run: install
	cd ./examples/app; pkger

tidy:
	$(GO_BIN) mod tidy -v

build: tidy
	cd ./cmd/pkger && $(GO_BIN) build $(GO_BUILD_FLAGS) -v .
	make tidy

test: tidy
	$(GO_BIN) test -count 1 -cover -tags ${TAGS} -timeout 1m ./...
	make tidy

cov:
	$(GO_BIN) test -coverprofile cover.out -tags ${TAGS} ./...
	go tool cover -html cover.out
	make tidy

ci-test:
	$(GO_BIN) test -tags ${TAGS} -race ./...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --enable-all
	make tidy

update:
	rm go.*
	$(GO_BIN) mod init
	$(GO_BIN) mod tidy
	make test
	make install
	make tidy

release-test:
	$(GO_BIN) test -tags ${TAGS} -race ./...
	make tidy

release:
	$(GO_BIN) get github.com/gobuffalo/release
	make tidy
	release -y -f version.go --skip-packr
	make tidy


