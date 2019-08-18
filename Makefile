TAGS ?= ""
GO_BIN ?= "go"

install: tidy
	cd ./cmd/pkger && $(GO_BIN) install -tags ${TAGS} -v .
	make tidy

tidy:
	$(GO_BIN) mod tidy

build: tidy
	$(GO_BIN) build -v .
	make tidy

test: tidy
	$(GO_BIN) test -cover -tags ${TAGS} -timeout 5s ./...
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



